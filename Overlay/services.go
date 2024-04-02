package overlay

import (
	"context"
	"errors"
	"log"

	pb "github.com/girivad/go-chord/Proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Lookup Services
// TO-DO: Update to handle timed contexts/SLAs?.

func (chordServer *ChordServer) FindSuccessor(ctx context.Context, keyHash *pb.Hash) (*pb.IP, error) {
	// Find the nearest predecessor and return its successor.
	log.Printf("[DEBUG] Find Successor Invoked")
	defer log.Printf("[DEBUG] Find Successor Completed.")

	// Ask the latest finger before the key to find the successor.
	for finger := chordServer.Capacity - 1; finger >= 0; finger-- {
		chordServer.FingerMuxs[finger].RLock()

		if chordServer.FingerTable[finger] == nil || chordServer.FingerTable[finger].Ip == chordServer.IP {
			chordServer.FingerMuxs[finger].RUnlock()
			continue
		}

		if isBetween(hash(chordServer.FingerTable[finger].Ip, chordServer.Capacity), chordServer.Hash, keyHash.Hash.Value) {
			log.Printf("[INFO] Find Successor Transferred to %s", chordServer.FingerTable[finger].Ip)
			ipMsg, err := chordServer.FingerTable[finger].LookupClient.FindSuccessor(ctx, keyHash)
			return ipMsg, err
		}

		chordServer.FingerMuxs[finger].RUnlock()
	}

	chordServer.FingerMuxs[0].RLock()
	successorIP := chordServer.FingerTable[0].Ip
	chordServer.FingerMuxs[0].RUnlock()

	// If the key is between me and my successor, return my successor.
	return &pb.IP{Ip: &wrapperspb.StringValue{Value: successorIP}}, nil
}

// Predecessor Services

func (chordServer *ChordServer) GetPredecessor(ctx context.Context, empty *emptypb.Empty) (*pb.IP, error) {
	log.Printf("[DEBUG] Get Predecessor Invoked.")
	defer log.Printf("[DEBUG] Get Predecessor Completed.")

	chordServer.PredecessorMux.RLock()
	if chordServer.Predecessor != nil {
		predecessorIP := chordServer.Predecessor.Ip
		chordServer.PredecessorMux.RUnlock()

		return &pb.IP{
			Ip: &wrapperspb.StringValue{Value: predecessorIP},
		}, nil
	}

	chordServer.PredecessorMux.RUnlock()

	return nil, errors.New("predecessor not known")
}

func (chordServer *ChordServer) UpdatePredecessor(ctx context.Context, ip *pb.IP) (*emptypb.Empty, error) {

	log.Printf("[DEBUG] Update Predecessor Invoked.")
	defer log.Printf("[DEBUG] Update Predecessor Completed.")

	var predecessorIP string

	chordServer.PredecessorMux.RLock()

	if chordServer.Predecessor != nil {
		predecessorIP = chordServer.Predecessor.Ip
	} else {
		predecessorIP = ""
	}

	chordServer.PredecessorMux.RUnlock()

	if predecessorIP == "" || isBetween(hash(ip.Ip.Value, chordServer.Capacity), hash(predecessorIP, chordServer.Capacity), chordServer.Hash) {
		var err error
		newPredecessor, err := Connect(ip.Ip.Value)

		if err != nil {
			return &emptypb.Empty{}, err
		}

		data, err := chordServer.DataToTransfer(hash(newPredecessor.Ip, chordServer.Capacity))

		if err != nil {
			return &emptypb.Empty{}, err
		}

		newPredecessor.DataClient.TransferData(context.Background(), data)

		chordServer.PredecessorMux.Lock()
		chordServer.Predecessor = newPredecessor
		chordServer.PredecessorMux.Unlock()
	}

	return &emptypb.Empty{}, nil
}

// Check Service (check if node is still alive).

func (chordServer *ChordServer) LiveCheck(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// Data Service: Transfer data to new owner
func (chordServer *ChordServer) DataToTransfer(nodeHash int64) (*pb.KVMap, error) {
	// TO-DO: Delete the keys.
	// TO-DO^2: Have a different route to delete the keys or do so when acknowledged
	// If predecessor is nil, it should transfer

	var predHash int64

	chordServer.PredecessorMux.RLock()

	if chordServer.Predecessor != nil {
		predecessorIP := chordServer.Predecessor.Ip
		predHash = hash(predecessorIP, chordServer.Capacity)
	} else {
		predHash = (hash(chordServer.IP, chordServer.Capacity) + 1) % (1 << chordServer.Capacity)
	}

	chordServer.PredecessorMux.RUnlock()

	keys := chordServer.keyIndex.KeysToTransfer(predHash, nodeHash, chordServer.Hash)

	keyValuePairs, err := chordServer.KVStore.GetValuesForTransfer(keys)

	return keyValuePairs, err
}

func (chordServer *ChordServer) TransferData(ctx context.Context, data *pb.KVMap) (*emptypb.Empty, error) {
	err := chordServer.KVStore.PutValuesForTransfer(data)
	chordServer.keyIndex.InsertBatch(data, func(key string) int64 { return hash(key, chordServer.Capacity) })

	return &emptypb.Empty{}, err
}
