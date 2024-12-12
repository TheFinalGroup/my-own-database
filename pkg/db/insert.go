package db

import (
	"bytes"

	btree "github.com/TheFinalGroup/my-own-database/pkg/b-tree"
)

func TreeInsert(node btree.BNode, key []byte, val []byte) btree.BNode {
	index := FindLessEqualNode(node, key)
	return InsertIntoNode(node, key, val, index)
}

func InsertIntoNode(node btree.BNode, key []byte, val []byte, index uint16) btree.BNode {
	newNode := btree.BNode{data: make([]byte, 2*BtreePageSize)}

	switch node.btype() {
	case BNodeLeaf:
		if bytes.Equal(key, node.getKey(index)) {
			leafNodeUpdate(&newNode, node, index, key, val)
		} else {
			leafNodeInsert(&newNode, node, index+1, key, val)
		}
	case BNodeInternal:
		internalNodeInsert(&newNode, node, index, key, val)
	default:
		panic("bad node!")
	}

	return newNode
}

func LeafNodeInsert(new *btree.BNode, old btree.BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNodeLeaf, old.numberOfKeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.numberOfKeys()-idx)
}

func LeafNodeUpdate(new *btree.BNode, old btree.BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNodeLeaf, old.numberOfKeys())
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx+1, old.numberOfKeys()-idx-1)
}

func InternalNodeInsert(new *btree.BNode, old btree.BNode, idx uint16, key []byte, val []byte) {
	childPointer := old.getPointer(idx)
	childNode := tree.get(childPointer)
	tree.del(childPointer)
	childNode = treeInsert(childNode, key, val)
	numberOfSplits, splitedNodes := splitNode(childNode)
	nodeReplaceChildren(new, old, idx, splitedNodes[:numberOfSplits]...)
}
