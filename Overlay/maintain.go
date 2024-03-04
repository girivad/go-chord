package overlay

import (
	"context"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const period time.Duration = 10 * time.Second
const MaxRetries int = 3

// Implement "Notify" (notifies a node that the caller thinks it is their predecessor)

func (chordServer *ChordServer) Notify() {
	for {
		time.Sleep(period)

		_, err := chordServer.FingerTable[0].PredecessorClient.UpdatePredecessor(context.Background(), &wrapperspb.StringValue{
			Value: chordServer.IP,
		})

		if err != nil {
			log.Printf("[INFO] %s unable to notify succesor %s due to err %v", chordServer.IP, chordServer.FingerTable[0].Ip, err)
		} else {
			log.Printf("[INFO] %s notified %s that %s is its predecessor.", chordServer.IP, chordServer.FingerTable[0].Ip, chordServer.IP)
		}
	}
}

// Fix Finger Table (Periodically uses findSuccessor to update each finger)

func (chordServer *ChordServer) FixFingers() {
	retries := 0
	expired := false
	fingerToUpdate := 0

	isExpired := func() bool {
		retries++
		if retries > MaxRetries {
			retries = 0
			return true
		}

		return false
	}

	for {
		time.Sleep(period)

		fingerStart := (chordServer.Hash + 1<<(fingerToUpdate)) % (1 << chordServer.Capacity)

		if fingerToUpdate > 0 && !isBetween(hash(chordServer.FingerTable[fingerToUpdate-1].Ip), chordServer.Hash, fingerStart) {
			newFinger, err := Connect(chordServer.FingerTable[fingerToUpdate-1].Ip)
			if err == nil {
				retries = 0
				chordServer.FingerTable[fingerToUpdate] = newFinger
				log.Printf("[INFO] %s copied the previous finger %s to finger %d.", chordServer.IP, newFinger.Ip, fingerToUpdate)
				return
			}

			expired = isExpired()
			if expired {
				log.Printf("[INFO] Retries for %dth finger expired, moving to next finger.", fingerToUpdate)
				fingerToUpdate = (fingerToUpdate + 1)
			}
			continue
		}

		newFingerIp, err := chordServer.FindSuccessor(context.Background(), &wrapperspb.Int64Value{Value: fingerStart})

		if err != nil {
			log.Printf("[INFO] %s Unable to find %d finger due to %v, retrying...", chordServer.IP, fingerToUpdate, err)
		}

		newFinger, err := Connect(newFingerIp.Value)

		if err != nil {
			log.Printf("[INFO] %s Unable to connect to found %d finger %s due to %v, retrying...", chordServer.IP, fingerToUpdate, newFingerIp.Value, err)

			expired = isExpired()
			if expired {
				log.Printf("[INFO] Retries for %dth finger expired, moving to next finger.", fingerToUpdate)
				fingerToUpdate = (fingerToUpdate + 1)
			}
			continue
		}

		chordServer.FingerTable[fingerToUpdate] = newFinger

		log.Printf("[INFO] %s updated finger %d to %s", chordServer.IP, fingerToUpdate, newFinger.Ip)

		fingerToUpdate = (fingerToUpdate + 1) % (1 << chordServer.Capacity)

		log.Printf("[INFO] %s to update finger %d next", chordServer.IP, fingerToUpdate)
	}
}

// Check Predecessor (Set predecessor to nil if it is not live any more)

func (chordServer *ChordServer) CheckPredecessor() {
	retries := 0

	isExpired := func() bool {
		retries++

		if retries > MaxRetries {
			retries = 0
			return true
		}
		return false
	}

	for {
		time.Sleep(period)

		if chordServer.Predecessor == nil {
			continue
		}

		_, err := chordServer.Predecessor.CheckClient.LiveCheck(context.Background(), &emptypb.Empty{})

		if err != nil {

			if !isExpired() {
				log.Printf("[INFO] %s's predecessor did not respond to liveness check due to %v, retrying...", chordServer.IP, err)
				continue
			}

			log.Printf("[INFO] %s's predecessor did not respond to liveness check due to %v and was set to nil", chordServer.IP, err)
			chordServer.Predecessor = nil
			retries = 0
			continue
		}

		log.Printf("[INFO] %s's predecessor %s is still live.", chordServer.IP, chordServer.Predecessor.Ip)
	}
}

// Stabilize (Get successor's predecessor and set as my own - stabilizes after join in between)
func (chordServer *ChordServer) Stabilize() {
	for {
		time.Sleep(period)

		newSuccessorIp, err := chordServer.FingerTable[0].PredecessorClient.GetPredecessor(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Printf("[INFO] %s's successor %s failed to provide its predecessor due to %v", chordServer.IP, chordServer.FingerTable[0].Ip, err)
			continue
		}

		if !isBetween(hash(newSuccessorIp.Value), chordServer.Hash, hash(chordServer.FingerTable[0].Ip)) {
			// chordServer is still the latest predecessor to its successor (i.e. no new nodes have joined in between them).
			log.Printf("[INFO] %s is still the latest predecessor to %s.", chordServer.IP, chordServer.FingerTable[0].Ip)
			continue
		}

		newSuccessor, err := Connect(newSuccessorIp.Value)

		if err != nil {
			log.Printf("[INFO] %s failed to connect with its new successor %s, retrying while retaining the old successor...", chordServer.IP, newSuccessorIp.Value)
		}

		chordServer.FingerTable[0] = newSuccessor

		log.Printf("[INFO] %s is %s's new successor.", chordServer.FingerTable[0].Ip, chordServer.IP)
	}
}