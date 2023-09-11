package storage

import (
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
