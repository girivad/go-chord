package overlay

type BST struct {
	Key    string
	Hash   int64
	Left   *BST
	Right  *BST
	Parent *BST
}

func NewKeyIndex() *BST {
	return &BST{}
}

func (bst *BST) Insert(key string) {
	// Insert key into BST
}

func (bst *BST) Search(newHash int64, succHash int64) []string {
	keys := bst.SearchRight(newHash)
	if succHash < newHash {
		keys = append(keys, bst.SearchLeft(succHash)...)
	}

	return keys
}

func (bst *BST) SearchRight(newHash int64) []string {

	var keys []string

	// If it is less than the newHash, return searchRight of the right subtree.
	if bst.Hash < newHash {
		return bst.Right.SearchRight(newHash)
	}

	// If it is equal to the newHash, include the key.
	if bst.Hash == newHash {
		keys = append(keys, bst.Key)
	}

	// Return the traversal of the right subtree.
	if bst.Right != nil {
		keys = append(keys, bst.Right.Traverse()...)
	}

	return keys
}

func (bst *BST) SearchLeft(succHash int64) []string {
	var keys []string

	return keys
}

func (bst *BST) Traverse() []string {
	var traversal []string

	traversal = append(traversal, bst.Key)

	if bst.Left != nil {
		traversal = append(traversal, bst.Left.Traverse()...)
	}

	if bst.Right != nil {
		traversal = append(traversal, bst.Right.Traverse()...)
	}

	return traversal
}

func (bst *BST) Delete(key string) {

}
