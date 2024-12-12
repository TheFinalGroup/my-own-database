package btree

import (
	"encoding/binary"

	"github.com/TheFinalGroup/my-own-database/pkg/assertions"
	"github.com/TheFinalGroup/my-own-database/pkg/utils"
)

type BNode struct {
	data []byte
}

func (node BNode) Btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

func (node BNode) Nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) SetHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

func (node BNode) GetPtr(idx uint16) uint64 {
	assertions.Assert(idx < node.Nkeys())
	pos := utils.HEADER + 8*idx

	return binary.LittleEndian.Uint64(node.data[pos:])
}

func (node BNode) SetPtr(idx uint16, val uint64) {
	assertions.Assert(idx < node.Nkeys())
	pos := utils.HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

func OffsetPos(node BNode, idx uint16) uint16 {
	assertions.Assert(1 <= idx && idx <= node.Nkeys())
	return utils.HEADER + 8*node.Nkeys() + 2*(idx-1)
}

func (node BNode) GetOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}

	return binary.LittleEndian.Uint16(node.data[OffsetPos(node, idx):])
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[OffsetPos(node, idx):], offset)
}

func (node BNode) kvPos(idx uint16) uint16 {
	assertions.Assert(idx <= node.Nkeys())

	return utils.HEADER + 8*node.Nkeys() + 2*node.Nkeys() + node.GetOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	assertions.Assert(idx < node.Nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])

	return node.data[pos+4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	assertions.Assert(idx < node.Nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])

	return node.data[pos+4+klen:][:vlen]
}

func (node BNode) nbytes() uint16 {
	return node.kvPos(node.Nkeys())
}
