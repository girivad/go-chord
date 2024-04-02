package overlay

import (
	"log"
	"slices"
	"sync"

	pb "github.com/girivad/go-chord/Proto"
)

// To-do: Implement KeyIndex containing BST + Lock

type KeyIndex struct {
	root *BST
	lock sync.RWMutex
}

func NewKeyIndex() *KeyIndex {
	return &KeyIndex{
		root: new(BST),
	}
}

func (keyIndex *KeyIndex) Insert(key string, hash int64) {
	keyIndex.lock.Lock()
	keyIndex.root.Insert(key, hash, nil)
	keyIndex.lock.Unlock()
}

func (keyIndex *KeyIndex) Delete(key string, hash int64) {
	keyIndex.lock.Lock()
	keyIndex.root.Delete(key, hash)
	keyIndex.lock.Unlock()
}

func (keyIndex *KeyIndex) Visualize() {
	keyIndex.lock.RLock()
	keyIndex.root.Visualize()
	keyIndex.lock.RUnlock()
}

func (keyIndex *KeyIndex) KeysToTransfer(predHash, nodeHash, currHash int64) []string {
	keyIndex.lock.RLock()
	keys := keyIndex.root.KeysToTransfer(predHash, nodeHash, currHash)
	keyIndex.lock.RUnlock()
	return keys
}

func (keyIndex *KeyIndex) InsertBatch(data *pb.KVMap, hashFunc func(key string) int64) {
	var keys []string

	for key := range data.Kvmap {
		keys = append(keys, key)
	}

	keyIndex.lock.Lock()
	keyIndex.root.InsertBatch(keys, hashFunc)
	keyIndex.lock.Unlock()
}

type BST struct {
	Key    string
	Hash   int64
	Left   *BST
	Right  *BST
	Parent *BST
	Set    bool
}

func NewBST(key string, hash int64, left *BST, right *BST, parent *BST) *BST {
	keyIndex := &BST{Key: key, Hash: hash, Left: left, Right: right, Parent: parent, Set: true}
	return keyIndex
}

// SINGLE-KEY OPERATIONS: Insert Key, Delete Key

func (bst *BST) Insert(key string, hash int64, parent *BST) {
	if !bst.Set {
		bst.Key = key
		bst.Hash = hash
		bst.Parent = parent
		bst.Set = true
		return
	}

	if hash < bst.Hash {
		if bst.Left != nil {
			bst.Left.Insert(key, hash, bst)
		} else {
			bst.Left = NewBST(key, hash, nil, nil, bst)
		}
	} else if hash >= bst.Hash {
		if bst.Right != nil {
			bst.Right.Insert(key, hash, bst)
		} else {
			bst.Right = NewBST(key, hash, nil, nil, bst)
		}
	}
}

func (bst *BST) Leftmost() *BST {
	if bst.Left != nil && bst.Left.Set {
		return bst.Left.Leftmost()
	}

	return bst
}

func (bst *BST) Rightmost() *BST {
	if bst.Right != nil && bst.Right.Set {
		return bst.Right.Rightmost()
	}

	return bst
}

func (bst *BST) DeleteLeaf() {
	if bst.Parent == nil {
		bst.Set = false
		return
	}

	if bst.Parent.Hash < bst.Hash {
		bst.Parent.Right = nil
	} else {
		bst.Parent.Left = nil
	}
}

func (bst *BST) Delete(key string, hash int64) bool {
	var found bool
	var swapNode *BST

	if !bst.Set {
		return false
	}

	if bst.Hash < hash {
		if bst.Right != nil {
			found = bst.Right.Delete(key, hash)
		}
		return found
	}

	if bst.Hash > hash {
		if bst.Left != nil {
			found = bst.Left.Delete(key, hash)
		}
		return found
	}

	if bst.Key != key {
		if bst.Left != nil {
			found = bst.Left.Delete(key, hash)
		}
		if bst.Right != nil {
			found = found || bst.Right.Delete(key, hash)
		}
		return found
	}

	// Replace bst with the in-order successor/predecessor if bst.Key == key

	if bst.Right != nil && bst.Right.Set {
		swapNode = bst.Right.Leftmost()
	} else if bst.Left != nil && bst.Left.Set {
		swapNode = bst.Left.Rightmost()
	} else {
		swapNode = bst
	}

	bst.Key = swapNode.Key
	bst.Hash = swapNode.Hash
	swapNode.DeleteLeaf()

	return true
}

// BATCH-OPERATIONS: KeysToTransfer, (TO-DO) Insert Keys, Delete Keys

// Retrieve all keys between my predecessor and a new node
func (bst *BST) KeysToTransfer(predHash int64, newHash int64, currHash int64) []string {
	if currHash > predHash {
		return bst.KeysGreaterThan(newHash)
	}

	if newHash > predHash {
		return bst.KeysBetween(predHash, newHash)
	}

	return append(bst.KeysGreaterThan(predHash), bst.KeysLessThan(newHash)...)
}

func (bst *BST) KeysGreaterThan(lowerBound int64) []string {

	var keys []string

	if !bst.Set {
		return keys
	}

	// If it is less than or equal to the newHash, return searchRight of the right subtree.
	if bst.Hash <= lowerBound && bst.Right != nil {
		return bst.Right.KeysGreaterThan(lowerBound)
	} else if bst.Hash <= lowerBound {
		return keys
	}

	keys = append(keys, bst.Key)

	// Return the traversal of the right subtree, and any elements of the left subtree that might be greater than the lower bound.
	if bst.Right != nil {
		keys = append(keys, bst.Right.AllKeys()...)
	}
	if bst.Left != nil {
		keys = append(keys, bst.Left.KeysGreaterThan(lowerBound)...)
	}

	return keys
}

func (bst *BST) KeysLessThan(upperBound int64) []string {
	var keys []string

	if !bst.Set {
		return keys
	}

	if bst.Hash > upperBound && bst.Left != nil {
		return bst.Left.KeysLessThan(upperBound)
	} else if bst.Hash > upperBound {
		return keys
	}

	keys = append(keys, bst.Key)
	if bst.Right != nil {
		keys = append(keys, bst.Right.KeysLessThan(upperBound)...)
	}
	if bst.Left != nil {
		keys = append(keys, bst.Left.AllKeys()...)
	}

	return keys
}

// Finds all keys in (startHash, endHash]
func (bst *BST) KeysBetween(startHash int64, endHash int64) []string {
	var keys []string

	if !bst.Set {
		return keys
	}

	// If Key too small:
	if bst.Hash <= startHash && bst.Right != nil {
		return bst.Right.KeysBetween(startHash, endHash)
	} else if bst.Hash < startHash {
		return keys
	}

	// If Key too large:
	if bst.Hash > endHash && bst.Left != nil {
		return bst.Left.KeysBetween(startHash, endHash)
	} else if bst.Hash > endHash {
		return keys
	}

	// If Key in range:

	keys = append(keys, bst.Key)

	if bst.Right != nil {
		keys = append(keys, bst.Right.KeysLessThan(endHash)...)
	}
	if bst.Left != nil {
		keys = append(keys, bst.Left.KeysGreaterThan(startHash)...)
	}

	return keys
}

func (bst *BST) AllKeys() []string {
	var traversal []string
	if !bst.Set {
		return traversal
	}

	traversal = append(traversal, bst.Key)

	if bst.Left != nil {
		traversal = append(traversal, bst.Left.AllKeys()...)
	}

	if bst.Right != nil {
		traversal = append(traversal, bst.Right.AllKeys()...)
	}

	return traversal
}

func (bst *BST) Visualize() {
	if !bst.Set {
		return
	}

	var leftKey, rightKey string
	if bst.Left != nil && bst.Left.Set {
		leftKey = bst.Left.Key
	}
	if bst.Right != nil && bst.Right.Set {
		rightKey = bst.Right.Key
	}

	log.Printf("[INFO] %s: L %s R %s", bst.Key, leftKey, rightKey)

	if bst.Left != nil {
		bst.Left.Visualize()
	}
	if bst.Right != nil {
		bst.Right.Visualize()
	}
}

func (bst *BST) InsertBatch(keys []string, hashFunc func(key string) int64) {
	if len(keys) == 0 {
		return
	}

	if !bst.Set {
		slices.SortFunc(keys, func(key1, key2 string) int {
			if hashFunc(key1) < hashFunc(key2) {
				return -1
			}
			if hashFunc(key1) == hashFunc(key2) {
				return 0
			}
			return 1
		})

		midKey := keys[len(keys)/2]
		bst.Insert(midKey, hashFunc(midKey), bst.Parent)
	}

	var leftKeys, rightKeys []string
	var leftNodeKey, rightNodeKey string

	for _, key := range keys {
		if hashFunc(key) < bst.Hash {
			leftKeys = append(leftKeys, key)
			if key > leftNodeKey {
				leftNodeKey = key
			}
		} else if hashFunc(key) > bst.Hash {
			rightKeys = append(rightKeys, key)
			if key < rightNodeKey {
				rightNodeKey = key
			}
		}
	}

	if bst.Left == nil {
		bst.Left = NewBST(leftNodeKey, hashFunc(leftNodeKey), nil, nil, bst)
	}

	bst.Left.InsertBatch(leftKeys, hashFunc)

	if bst.Right == nil {
		bst.Right = NewBST(rightNodeKey, hashFunc(rightNodeKey), nil, nil, bst)
	}

	bst.Right.InsertBatch(rightKeys, hashFunc)
}
