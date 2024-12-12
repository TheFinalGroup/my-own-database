package btree

import (
	"bytes"
	"encoding/binary"

	"github.com/TheFinalGroup/my-own-database/pkg/assertions"
	"github.com/TheFinalGroup/my-own-database/pkg/utils"
)

type BTree struct {
	root uint64
	get  func(uint64) BNode // dereference a pointer
	new  func(BNode) uint64 // allocate a new page
	del  func(uint64)       // deallocate a page
}

func init() {
	node1max := utils.HEADER + 8 + 2 + 4 + utils.BTREE_MAX_KEY_SIZE + utils.BTREE_MAX_VAL_SIZE
	assertions.Assert(node1max <= utils.BTREE_PAGE_SIZE)
}

func FindLessEqualNode(node BNode, key []byte) uint16 {
	nkeys := node.Nkeys()
	found := uint16(0)
	// the first key is a copy from the parent node,
	// thus it's always less than or equal to the key.
	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}
	return found
}

func BisectFindIndexLessEqual(node BNode, key []byte) uint16 {
	start, end := uint16(0), node.Nkeys()-1
	for start <= end {
		mid := start + (end-start)/2

		compare := bytes.Compare(node.getKey(mid), key)

		if compare == 0 {
			return mid
		}

		if compare > 0 {
			end = mid - 1
		} else {
			if end == start {
				return mid
			}
			start = mid + 1
		}
	}
	panic("no index found")
}

func SplitNode(node BNode) (uint16, [3]BNode) {
	if node.nbytes() <= BtreePageSize {
		node.data = node.data[:BtreePageSize]
		return 1, [3]BNode{node}
	}

	left := BNode{make([]byte, 2*BtreePageSize)}
	right := BNode{make([]byte, BtreePageSize)}
	splitIntoTwoNodes(&left, &right, node)
	if left.nbytes() <= BtreePageSize {
		left.data = left.data[:BtreePageSize]
		return 2, [3]BNode{left, right}
	}

	leftLeft := BNode{make([]byte, BtreePageSize)}
	middle := BNode{make([]byte, BtreePageSize)}
	splitIntoTwoNodes(&leftLeft, &middle, left)
	assertions.Assert(leftLeft.nbytes() <= BtreePageSize)

	return 3, [3]BNode{leftLeft, middle, right}
}

func SplitIntoTwoNodes(left *BNode, right *BNode, old BNode) {
	numberOfKeys := GetSplitKeyIndex(old)

	right.setHeader(old.getType(), numberOfKeys)
	NodeAppendRange(right, old, 0, 0, numberOfKeys)

	left.setHeader(old.getType(), old.Nkeys()-numberOfKeys)
	NodeAppendRange(left, old, 0, numberOfKeys, old.Nkeys()-numberOfKeys)
}

func GetSplitKeyIndex(node BNode) uint16 {
	for i := uint16(1); i <= node.Nkeys(); i++ {
		if node.kvPos(i) > BtreePageSize {
			return i - 1
		}
	}
	return node.Nkeys()
}

func BisectFindSplitKeyIndex(node BNode) uint16 {
	start, end := uint16(0), node.Nkeys()-1
	for start <= end {
		mid := start + (end-start)/2

		pageSize := node.kvPos(mid)

		if pageSize == BtreePageSize {
			return mid
		}

		if pageSize > BtreePageSize {
			end = mid - 1
		} else {
			if end == start {
				return mid
			}
			start = mid + 1
		}
	}
	panic("no index found")
}

func NodeAppendRange(new *BNode, old BNode, newIndex uint16, oldIndex uint16, length uint16) {
	assertions.Assert(oldIndex+length <= old.Nkeys())
	assertions.Assert(newIndex+length <= new.Nkeys())
	if length == 0 {
		return
	}

	for i := uint16(0); i < length; i++ {
		new.setPtr(newIndex+i, old.getPtr(oldIndex+i))
	}

	dstBegin := new.getOffset(newIndex)
	srcBegin := old.getOffset(oldIndex)
	for i := uint16(1); i <= length; i++ { // NOTE: the range is [1, length]
		offset := dstBegin + old.getOffset(oldIndex+i) - srcBegin
		new.setOffset(newIndex+i, offset)
	}

	begin := old.kvPos(oldIndex)
	end := old.kvPos(oldIndex + length)
	copy(new.data[new.kvPos(newIndex):], old.data[begin:end])
}

func nodeAppendKV(new *BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	new.setPtr(idx, ptr)
	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new.data[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new.data[pos+2:], uint16(len(val)))
	copy(new.data[pos+4:], key)
	copy(new.data[pos+4+uint16(len(key)):], val)
	new.setOffset(idx+1, new.getOffset(idx)+4+uint16(len(key)+len(val)))
}

func NodeReplaceChildren(new *BNode, old BNode, idx uint16, children ...BNode) {
	numberOfChildren := uint16(len(children))
	new.setHeader(BNodeInternal, old.Nkeys()+numberOfChildren-1)
	NodeAppendRange(new, old, 0, 0, idx)
	for i, child := range children {
		nodeAppendKV(new, idx+uint16(i), tree.new(child), child.getKey(0), nil)
	}
	NodeAppendRange(new, old, idx+numberOfChildren, idx+1, old.Nkeys()-(idx+1))
}
