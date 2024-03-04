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
		if isBetween(hash(chordServer.FingerTable[finger].Ip), chordServer.Hash, keyHash.Value) {
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
	if isBetween(hash(ip.Value), hash(chordServer.Predecessor.Ip), chordServer.Hash) {
		var err error
		chordServer.Predecessor, err = Connect(ip.Value)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

// Check Service (check if node is still alive).

func (chordServer *ChordServer) LiveCheck(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// Data Service: Transfer data to new owner
func (chordServer *ChordServer) TransferData(ctx context.Context, nodeHash *wrapperspb.Int64Value) (*pb.KVMap, error) {
	return nil, nil
}
