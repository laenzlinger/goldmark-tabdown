package tabdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type tabdownParser struct {
}

var defaultTabdownParser = &tabdownParser{}

func (b *tabdownParser) Trigger() []byte {
	return []byte{'['}
}

func isChordLine(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	chord := false
	for i := 0; i < len(line); i++ {
		if line[i] == '[' {
			chord = true
		} else if line[i] == ']' {
			chord = false
		} else {
			if !chord && util.IsSpace(line[i]) {
				return false
			}
		}
	}
	return true
}

func (b *tabdownParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	linenum, _ := reader.Position()
	if linenum != 0 {
		return nil, parser.NoChildren
	}
	line, _ := reader.PeekLine()
	if isChordLine(line) {
		return gast.NewTextBlock(), parser.NoChildren
	}
	return nil, parser.NoChildren
}

func (b *tabdownParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if !util.IsBlank(line) {
		reader.Advance(segment.Len())
		return parser.Close
	}
	node.Lines().Append(segment)
	return parser.Continue | parser.NoChildren
}

func (b *tabdownParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
	lines := node.Lines()
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}
}

func (b *tabdownParser) CanInterruptParagraph() bool {
	return true
}

func (b *tabdownParser) CanAcceptIndentedLine() bool {
	return true
}

// NewParser returns a BlockParser that can parse Tabdown blocks.
func NewParser() parser.BlockParser {
	return defaultTabdownParser
}

type tabdown struct {
}

// Option is a functional option type for this extension.
type Option func(*tabdown)

// Tabdown is an extension for goldmark.
var Tabdown = &tabdown{}

// New returns a new tabdown extension.
func New(opts ...Option) goldmark.Extender {
	e := &tabdown{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *tabdown) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewParser(), 0),
		),
	)
}
