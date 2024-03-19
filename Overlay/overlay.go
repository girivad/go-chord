package overlay

import (
	"context"
	"fmt"
	"net"

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
	Ip                string
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
	KeyIndex    *BST
	pb.UnimplementedLookupServer
	pb.UnimplementedPredecessorServer
	pb.UnimplementedCheckServer
	pb.UnimplementedDataServer
}

func (chordServer *ChordServer) RegisterKey(key string) {
	if !isBetween(hash(key, chordServer.Capacity), hash(chordServer.Predecessor.Ip, chordServer.Capacity), hash(chordServer.IP, chordServer.Capacity)) {
		fmt.Printf("Attempted to register Key %s with node %s, but doesn't belong here.\n", key, chordServer.IP)
		return
	}

	chordServer.KeyIndex.Insert(key)
}

func (chordServer *ChordServer) RegisterDelete(key string) {
	if !isBetween(hash(key, chordServer.Capacity), hash(chordServer.Predecessor.Ip, chordServer.Capacity), hash(chordServer.IP, chordServer.Capacity)) {
		fmt.Printf("Attempted to register delete key %s at node %s, but doesn't belong here.\n", key, chordServer.IP)
		return
	}

	chordServer.KeyIndex.Delete(key)
}

func NewChordServer(ip string, capacity int64) *ChordServer {
	chordServer := &ChordServer{
		IP:          ip,
		Hash:        hash(ip, capacity),
		Capacity:    capacity,
		Predecessor: nil,
		FingerTable: make([]*ChordNode, capacity),
		KeyIndex:    &BST{},
	}

	chordServer.KVStore = data.NewDataServer(chordServer.RegisterKey, chordServer.RegisterDelete)

	return chordServer
}

func (chordServer *ChordServer) Serve() error {
	// Serve the gRPC server as well, at port 8081.
	grpcListener, err := net.Listen("tcp", ":8081")

	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPredecessorServer(grpcServer, chordServer)
	pb.RegisterLookupServer(grpcServer, chordServer)
	pb.RegisterCheckServer(grpcServer, chordServer)
	pb.RegisterDataServer(grpcServer, chordServer)

	// Data served from port 8080.
	go chordServer.KVStore.Serve(8080)
	go chordServer.Notify()
	go chordServer.FixFingers()
	go chordServer.CheckPredecessor()
	go chordServer.Stabilize()

	err = grpcServer.Serve(grpcListener)

	return err
}

func Connect(ip string) (*ChordNode, error) {
	// Returns pointer to ChordNode with clients to the IP address.
	clientConn, err := grpc.Dial(ip+":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	chordNode := &ChordNode{
		Ip:                ip,
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
	chordServer.FingerTable[0], err = Connect(successorIpMsg.Value)

	if err != nil {
		return err
	}

	// Ask successor for its predecessor, set as own
	predecessorIpMsg, err := chordServer.FingerTable[0].PredecessorClient.GetPredecessor(context.Background(), &emptypb.Empty{})

	if err != nil {
		return err
	}

	chordServer.Predecessor, err = Connect(predecessorIpMsg.Value)

	return err
}
