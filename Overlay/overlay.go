package overlay

import (
	"context"

	data "github.com/girivad/go-chord/Data"
	pb "github.com/girivad/go-chord/Proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Interface for nodes in the Chord Ring.
type ChordNode struct {
	// Collection of clients for multiple services.
	PredecessorClient pb.PredecessorClient
	LookupClient      pb.LookupClient
	CheckClient       pb.CheckClient
	DataClient        pb.DataClient
}

// The local server
type ChordServer struct {
	KVStore     *data.DataServer
	IP          string
	Hash        int64
	Capacity    int64
	Predecessor *ChordNode
	FingerTable []*ChordNode
}

func hash(ip string) int64 {
	// Placeholder Hash
	return 0
}

func NewChordServer(ip string, capacity int64) *ChordServer {
	return &ChordServer{
		KVStore:     data.NewDataServer(),
		IP:          ip,
		Hash:        hash(ip),
		Capacity:    capacity,
		Predecessor: nil,
		FingerTable: make([]*ChordNode, capacity),
	}
}

func (chordServer *ChordServer) Serve() {
	// Data served from port 8080.
	chordServer.KVStore.Serve(8080)
	// Serve the gRPC server as well.
}

func (chordServer *ChordServer) Connect(ip string) (*ChordNode, error) {
	// Returns pointer to ChordNode with clients to the IP address.
	clientConn, err := grpc.Dial(ip+":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	chordNode := &ChordNode{
		PredecessorClient: pb.NewPredecessorClient(clientConn),
		LookupClient:      pb.NewLookupClient(clientConn),
		CheckClient:       pb.NewCheckClient(clientConn),
		DataClient:        pb.NewDataClient(clientConn),
	}

	return chordNode, nil
}

func (chordServer *ChordServer) Join(contactNode *ChordNode) error {
	// Find successor
	successorIpMsg, err := contactNode.LookupClient.FindSuccessor(context.Background(), &(wrapperspb.Int64Value{Value: chordServer.Hash}))

	if err != nil {
		return err
	}

	// Set successor
	chordServer.FingerTable[0], err = chordServer.Connect(successorIpMsg.Value)

	if err != nil {
		return err
	}

	// Ask successor for its predecessor, set as own
	predecessorIpMsg, err := chordServer.FingerTable[0].PredecessorClient.GetPredecessor(context.Background(), &emptypb.Empty{})

	if err != nil {
		return err
	}

	chordServer.Predecessor, err = chordServer.Connect(predecessorIpMsg.Value)

	return err
}