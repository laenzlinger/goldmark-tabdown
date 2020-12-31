package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

// A ChordBlock
type ChordBlock struct {
	gast.BaseBlock
}

// Dump implements Node.Dump .
func (n *ChordBlock) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindChordBlock is a NodeKind of the ChordBlock node.
var KindChordBlock = gast.NewNodeKind("ChordBlock")

// Kind implements Node.Kind.
func (n *ChordBlock) Kind() gast.NodeKind {
	return KindChordBlock
}

// NewChordBlock returns a new ChordBlock node.
func NewChordBlock() *ChordBlock {
	return &ChordBlock{
		BaseBlock: gast.BaseBlock{},
	}
}
