package overlay

import (
	"github.com/girivad/go-chord/data"
)

type ChordNode struct {
	KVStore       *data.DataServer
	IP            string
	Hash          int
	Capacity      int
	PredecessorIP string
	FingerTable   []string
}

func hash(ip string) int {
	// Placeholder Hash
	return 0
}

func NewChordNode(ip string, capacity int) *ChordNode {
	return &ChordNode{
		KVStore:       data.NewDataServer(),
		IP:            ip,
		Hash:          hash(ip),
		Capacity:      capacity,
		PredecessorIP: "",
		FingerTable:   make([]string, capacity),
	}
}
