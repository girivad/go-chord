package overlay

import "fmt"

type BST struct {
	Key    string
	Hash   int64
	Left   *BST
	Right  *BST
	Parent *BST
	Set    bool
}

func NewKeyIndex(key string, hash int64, left *BST, right *BST, parent *BST) *BST {
	keyIndex := &BST{Key: key, Hash: hash, Left: left, Right: right, Parent: parent, Set: true}
	return keyIndex
}

// SINGLE-KEY OPERATIONS: Insert Key, Delete Key

func (bst *BST) Insert(key string, hash int64, parent *BST) {
	if !bst.Set {
		bst.Key = key
		bst.Hash = hash
		// Not quite sure that we need the parent.
		bst.Parent = parent
		bst.Set = true
		return
	}

	if hash < bst.Hash {
		if bst.Left != nil {
			bst.Left.Insert(key, hash, bst)
		} else {
			bst.Left = NewKeyIndex(key, hash, nil, nil, bst)
		}
	} else if hash > bst.Hash {
		if bst.Right != nil {
			bst.Right.Insert(key, hash, bst)
		} else {
			bst.Right = NewKeyIndex(key, hash, nil, nil, bst)
		}
	} else {
		fmt.Printf("Tried to insert duplicate key %s into BST", key)
	}
}

func (bst *BST) Delete(key string, hash int64) {
	if !bst.Set {
		return
	}

	// Replace bst with the in-order successor if bst.Key == key
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
