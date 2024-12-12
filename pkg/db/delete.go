package db

import (
	"bytes"

	btree "github.com/TheFinalGroup/my-own-database/pkg/b-tree"
)

func TreeDelete(tree *btree.BTree, node btree.BNode, key []byte) btree.BNode {
	idx := findLessEqualNode(node, key)
	switch node.btype() {
	case BNodeLeaf:
		if !bytes.Equal(key, node.getKey(idx)) {
			return btree.BNode{}
		}
		newNode := btree.BNode{data: make([]byte, BtreePageSize)}
		leafDelete(newNode, node, idx)
		return newNode
	case BNodeInternal:
		return nodeDelete(tree, node, idx, key)
	default:
		panic("bad node!")
	}
}

func LeafDelete(new btree.BNode, old btree.BNode, idx uint16) {
	new.setHeader(BNodeLeaf, old.nkeys()-1)
	nodeAppendRange(&new, old, 0, 0, idx)
	nodeAppendRange(&new, old, idx, idx+1, old.numberOfKeys()-(idx+1))
}

func nodeDelete(bTree *btree.BTree, node btree.BNode, idx uint16, key []byte) btree.BNode {
	return btree.BNode{}
}
