package overlay

import (
	"context"
	"log"
	"time"

	pb "github.com/girivad/go-chord/Proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const period time.Duration = 10 * time.Second
const MaxRetries int = 3

// Implement "Notify" (notifies a node that the caller thinks it is their predecessor)

func (chordServer *ChordServer) Notify() {
	for {
		time.Sleep(period)
		log.Printf("[INFO] Notifying successor %s", chordServer.FingerTable[0].Ip)
		chordServer.FingerMuxs[0].RLock()
		successorIP := chordServer.FingerTable[0].Ip

		if successorIP == chordServer.IP {
			log.Printf("[INFO] No successor to notify.")
			chordServer.FingerMuxs[0].RUnlock()
			continue
		}

		_, err := chordServer.FingerTable[0].PredecessorClient.UpdatePredecessor(context.Background(), &pb.IP{
			Ip: &wrapperspb.StringValue{Value: chordServer.IP},
		})

		chordServer.FingerMuxs[0].RUnlock()

		if err != nil {
			log.Printf("[DEBUG] Unable to notify succesor %s due to err: %v", successorIP, err)
		} else {
			log.Printf("[INFO] Notified successor %s", successorIP)
		}
	}
}

// Fix Finger Table (Periodically uses findSuccessor to update each finger)

func (chordServer *ChordServer) FixFingers() {
	retries := 0
	expired := false

	var fingerToUpdate, fingerStart int64
	var newFinger *ChordNode
	var previousFingerIP string

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
		log.Printf("[INFO] Fixing Finger %d...", fingerToUpdate)

		fingerStart = (chordServer.Hash + 1<<(fingerToUpdate)) % (1 << chordServer.Capacity)

		if fingerToUpdate > 0 {
			chordServer.FingerMuxs[fingerToUpdate-1].RLock()
			previousFingerIP = chordServer.FingerTable[fingerToUpdate-1].Ip
		}

		if fingerToUpdate > 0 && isBetween(hash(previousFingerIP, chordServer.Capacity), fingerStart, chordServer.Hash) { // Use the previously updated finger if in the right segment of the ring.
			chordServer.FingerMuxs[fingerToUpdate-1].RUnlock()
			newFinger, err := Connect(previousFingerIP) // Could do chordServer.FingerTable[fingerToUpdate - 1], but the lock would tremendously slow down most operations if all fingers are the same node. Worth thinking about.
			if err == nil {
				chordServer.FingerMuxs[fingerToUpdate].Lock()
				chordServer.FingerTable[fingerToUpdate] = newFinger
				chordServer.FingerMuxs[fingerToUpdate].Unlock()

				log.Printf("[INFO] %s copied the previous finger %s to finger %d.", chordServer.IP, newFinger.Ip, fingerToUpdate)

				retries = 0
				continue
			}

			expired = isExpired()
			if expired {
				log.Printf("[INFO] Retries for %dth finger expired, moving to next finger.", fingerToUpdate)
				fingerToUpdate = (fingerToUpdate + 1) % chordServer.Capacity
			}
			continue
		} else if fingerToUpdate > 0 {
			chordServer.FingerMuxs[fingerToUpdate-1].RUnlock()
		}

		newFingerIp, err := chordServer.FindSuccessor(context.Background(), &pb.Hash{
			Hash: &wrapperspb.Int64Value{Value: fingerStart},
		})

		if err != nil {
			log.Printf("[INFO] %s Unable to find %d finger due to %v, retrying...", chordServer.IP, fingerToUpdate, err)
		}

		log.Printf("[INFO] FF: New Finger %d is %s", fingerToUpdate, newFingerIp.Ip.Value)

		if newFingerIp.Ip.Value == chordServer.IP {
			log.Printf("[INFO] %s still getting itself as finger %d", chordServer.IP, fingerToUpdate)
			continue
		}

		newFinger, err = Connect(newFingerIp.Ip.Value)
		if err != nil {
			log.Printf("[INFO] %s Unable to connect to found %d finger %s due to %v, retrying...", chordServer.IP, fingerToUpdate, newFingerIp.Ip.Value, err)

			expired = isExpired()
			if expired {
				log.Printf("[INFO] Retries for %dth finger expired, moving to next finger.", fingerToUpdate)
				fingerToUpdate = (fingerToUpdate + 1) % chordServer.Capacity
			}
			continue
		}
		chordServer.FingerMuxs[fingerToUpdate].Lock()
		chordServer.FingerTable[fingerToUpdate] = newFinger
		chordServer.FingerMuxs[fingerToUpdate].Unlock()

		log.Printf("[INFO] %s updated finger %d to %s", chordServer.IP, fingerToUpdate, newFinger.Ip)

		fingerToUpdate = (fingerToUpdate + 1) % chordServer.Capacity

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

		chordServer.PredecessorMux.RLock()

		if chordServer.Predecessor == nil {
			log.Printf("[INFO] No predecessor to check.")
			chordServer.PredecessorMux.RUnlock()
			continue
		}

		_, err := chordServer.Predecessor.CheckClient.LiveCheck(context.Background(), &emptypb.Empty{})

		chordServer.PredecessorMux.RUnlock()

		if err != nil {

			if !isExpired() {
				log.Printf("[INFO] %s's predecessor did not respond to liveness check due to %v, retrying...", chordServer.IP, err)
				continue
			}

			log.Printf("[INFO] %s's predecessor did not respond to liveness check due to %v and was set to nil", chordServer.IP, err)
			chordServer.PredecessorMux.Lock()
			chordServer.Predecessor = nil
			chordServer.PredecessorMux.Unlock()
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
		log.Println("[INFO] Stabilizing...")

		chordServer.FingerMuxs[0].RLock()
		successorIP := chordServer.FingerTable[0].Ip

		newSuccessorIp, err := chordServer.FingerTable[0].PredecessorClient.GetPredecessor(context.Background(), &emptypb.Empty{})

		chordServer.FingerMuxs[0].RUnlock()

		if err != nil {
			log.Printf("[INFO] %s's successor %s failed to provide its predecessor due to %v", chordServer.IP, successorIP, err)
			continue
		}

		if successorIP != chordServer.IP && !isBetween(hash(newSuccessorIp.Ip.Value, chordServer.Capacity), chordServer.Hash, hash(successorIP, chordServer.Capacity)) {
			// chordServer is still the latest predecessor to its successor (i.e. no new nodes have joined in between them).
			log.Printf("[INFO] %s is still the latest predecessor to %s.", chordServer.IP, successorIP)
			continue
		}

		newSuccessor, err := Connect(newSuccessorIp.Ip.Value)

		if err != nil {
			log.Printf("[INFO] %s failed to connect with its new successor %s, retrying while retaining the old successor...", chordServer.IP, newSuccessorIp.Ip.Value)
		}

		chordServer.FingerMuxs[0].Lock()

		chordServer.FingerTable[0] = newSuccessor

		chordServer.FingerMuxs[0].Unlock()

		log.Printf("[INFO] %s is %s's new successor.", newSuccessorIp, chordServer.IP)
	}
}
