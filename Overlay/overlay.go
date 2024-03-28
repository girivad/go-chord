package overlay

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	data "github.com/girivad/go-chord/Data"
	pb "github.com/girivad/go-chord/Proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	KVStore        *data.DataServer
	IP             string
	Hash           int64
	Capacity       int64
	Predecessor    *ChordNode
	FingerTable    []*ChordNode
	KeyIndex       *BST
	FingerMuxs     []*sync.RWMutex
	PredecessorMux *sync.RWMutex
	pb.UnimplementedLookupServer
	pb.UnimplementedPredecessorServer
	pb.UnimplementedCheckServer
	pb.UnimplementedDataServer
}

func NewChordServer(ip string, capacity int64) (*ChordServer, error) {
	chordServer := &ChordServer{
		IP:             ip,
		Hash:           hash(ip, capacity),
		Capacity:       capacity,
		Predecessor:    nil,
		FingerTable:    make([]*ChordNode, capacity),
		FingerMuxs:     make([]*sync.RWMutex, capacity),
		PredecessorMux: &sync.RWMutex{},
	}

	chordServer.KVStore = data.NewDataServer(chordServer.RegisterKey, chordServer.RegisterDelete)
	chordServer.KeyIndex = &BST{}

	successor, err := Connect(ip)
	if err != nil {
		return nil, err
	}
	chordServer.FingerTable[0] = successor

	for finger := 0; finger < int(capacity); finger++ {
		chordServer.FingerMuxs[finger] = &sync.RWMutex{}
	}

	return chordServer, nil
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
	successorIpMsg, err := contactNode.LookupClient.FindSuccessor(context.Background(), &pb.Hash{
		Hash: &(wrapperspb.Int64Value{Value: chordServer.Hash}),
	})

	if err != nil {
		return err
	}

	log.Printf("[INFO] %s joining chord ring of %s: successor is %s", chordServer.IP, contactNode.Ip, successorIpMsg.Ip.Value)

	// Set successor
	chordServer.FingerTable[0], err = Connect(successorIpMsg.Ip.Value)

	if err != nil {
		return err
	}

	return err
}

func (chordServer *ChordServer) Leave() {

	transferData := make(map[string]*pb.Value)

	// for key, value := range chordServer.KVStore.KVMap {
	// 	&pb.KVMap{}
	// }

	chordServer.FingerTable[0].DataClient.TransferData(context.Background(), &pb.KVMap{Kvmap: transferData})
}

func (chordServer *ChordServer) RegisterKey(key string) {
	if chordServer.Predecessor != nil && !isBetween(hash(key, chordServer.Capacity), hash(chordServer.Predecessor.Ip, chordServer.Capacity), hash(chordServer.IP, chordServer.Capacity)) {
		fmt.Printf("Attempted to register Key %s with node %s, but doesn't belong here.\n", key, chordServer.IP)
		return
	}

	chordServer.KeyIndex.Insert(key, hash(key, chordServer.Capacity), nil)
	log.Printf("[INFO] Post-Insert %s", key)
	chordServer.KeyIndex.Visualize()
}

func (chordServer *ChordServer) RegisterDelete(key string) {
	if chordServer.Predecessor != nil && !isBetween(hash(key, chordServer.Capacity), hash(chordServer.Predecessor.Ip, chordServer.Capacity), hash(chordServer.IP, chordServer.Capacity)) {
		fmt.Printf("Attempted to register delete key %s at node %s, but doesn't belong here.\n", key, chordServer.IP)
		return
	}

	chordServer.KeyIndex.Delete(key, hash(key, chordServer.Capacity))

	log.Printf("[INFO] Post-Delete %s", key)
	chordServer.KeyIndex.Visualize()
}
