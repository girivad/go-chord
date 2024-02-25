package overlay

import (
	"context"

	pb "github.com/girivad/go-chord/Proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Lookup Services

func (chordServer *ChordServer) FindSuccessor(ctx context.Context, hash *wrapperspb.Int64Value) (*wrapperspb.StringValue, error) {
	//
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
