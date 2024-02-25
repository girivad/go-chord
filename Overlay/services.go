package overlay

import (
	"context"

	pb "github.com/girivad/go-chord/Proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Lookup Services

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

}

func (chordServer *ChordServer) UpdatePredecessor(ctx context.Context, empty *wrapperspb.StringValue) (*emptypb.Empty, error) {

}

func (chordServer *ChordServer) LiveCheck(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {

}

func (chordServer *ChordServer) TransferData(ctx context.Context, hash *wrapperspb.Int64Value) (*pb.KVMap, error) {

}
