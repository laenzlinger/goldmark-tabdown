package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

// A Chord struct represents a chord with its lyrics
type Chord struct {
	gast.BaseBlock
	Name   []byte
	Indent int
}

// Dump implements Node.Dump.
func (n *Chord) Dump(source []byte, level int) {
	m := map[string]string{
		"Name":   fmt.Sprintf("%s", n.Name),
		"Indent": fmt.Sprintf("%v", n.Indent),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

func (n *Chord) EndsWithSpace(source []byte) bool {
	if n.LastChild() == nil {
		return false
	}
	if n.Lines().Len() <= 0 {
		return false
	}
	lastLine := n.Lines().At(n.Lines().Len() - 1)
	trimmedLine := lastLine.TrimRightSpace(source)
	return trimmedLine.Len() < lastLine.Len()
}

// KindChord is a NodeKind of the Chord node.
var KindChord = gast.NewNodeKind("Chord")

// Kind implements Node.Kind.
func (n *Chord) Kind() gast.NodeKind {
	return KindChord
}

// NewChord returns a new TaskCheckBox node.
func NewChord(indent int, name []byte) *Chord {
	return &Chord{
		BaseBlock: gast.BaseBlock{},
		Name:      name,
		Indent:    indent,
	}
}
