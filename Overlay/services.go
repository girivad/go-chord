package overlay

import (
	"context"
	"errors"

	pb "github.com/girivad/go-chord/Proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Lookup Services
// TO-DO: Update to handle timed contexts.

func (chordServer *ChordServer) FindSuccessor(ctx context.Context, keyHash *wrapperspb.Int64Value) (*wrapperspb.StringValue, error) {
	// Find the nearest predecessor and return its successor.

	// Ask the latest finger before the key to find the successor.
	for finger := chordServer.Capacity - 1; finger >= 0; finger-- {
		if isBetween(hash(chordServer.FingerTable[finger].Ip, chordServer.Capacity), chordServer.Hash, keyHash.Value) {
			ipMsg, err := chordServer.FingerTable[finger].LookupClient.FindSuccessor(ctx, keyHash)
			return ipMsg, err
		}
	}

	// If the key is between me and my successor, return my successor.
	return &wrapperspb.StringValue{Value: chordServer.FingerTable[0].Ip}, nil
}

// Predecessor Services

func (chordServer *ChordServer) GetPredecessor(ctx context.Context, empty *emptypb.Empty) (*wrapperspb.StringValue, error) {
	if chordServer.Predecessor != nil {
		return &wrapperspb.StringValue{
			Value: chordServer.Predecessor.Ip,
		}, nil
	}

	return nil, errors.New("predecessor not known")
}

func (chordServer *ChordServer) UpdatePredecessor(ctx context.Context, ip *wrapperspb.StringValue) (*emptypb.Empty, error) {
	if isBetween(hash(ip.Value, chordServer.Capacity), hash(chordServer.Predecessor.Ip, chordServer.Capacity), chordServer.Hash) {
		var err error
		newPredecessor, err := Connect(ip.Value)

		data, err := chordServer.DataToTransfer(hash(newPredecessor.Ip, chordServer.Capacity))

		newPredecessor.DataClient.TransferData(context.Background(), data)

		// Transfer Data to this new predecessor using routines used in TransferDataIn/Out below

		chordServer.Predecessor = newPredecessor

		return &emptypb.Empty{}, err
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

	predHash := hash(chordServer.Predecessor.Ip, chordServer.Capacity)
	keys := chordServer.KeyIndex.KeysToTransfer(predHash, nodeHash, chordServer.Hash)

	keyValuePairs, err := chordServer.KVStore.GetValuesForTransfer(keys)

	return keyValuePairs, err
}

func (chordServer *ChordServer) TransferData(ctx context.Context, data *pb.KVMap) (*emptypb.Empty, error) {
	err := chordServer.KVStore.PutValuesForTransfer(data)

	return &emptypb.Empty{}, err
}
