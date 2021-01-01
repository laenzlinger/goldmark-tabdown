package tabdown

import (
	"regexp"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/laenzlinger/goldmark-tabdown/ast"
)

type chordBlockParser struct {
}

var defaultChordBlockParser = &chordBlockParser{}

func (b *chordBlockParser) Trigger() []byte {
	return []byte{'['}
}

func isChordLine(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	chord := false
	nrOfChords := 0
	chordLength := 0
	for i := 0; i < len(line); i++ {
		if line[i] == '[' {
			chordLength = 0
			chord = true
		} else if chord && line[i] == ']' {
			if chordLength > 0 {
				nrOfChords++
			}
			chord = false
		} else {
			isSpace := util.IsSpace(line[i])
			if chord && !isSpace {
				chordLength++
			}
			if !chord && !isSpace {
				return false
			}
		}
	}
	return nrOfChords > 0
}

func (b *chordBlockParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if !isChordLine(line) {
		return nil, parser.NoChildren
	}

	reader.SkipSpaces()
	_, segment := reader.PeekLine()
	chordBlock := ast.NewChordBlock()
	chordBlock.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return chordBlock, parser.NoChildren
}

func (b *chordBlockParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if util.IsBlank(line) {
		reader.Advance(segment.Len() - 1)
		return parser.Continue | parser.NoChildren
	}
	node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return parser.Close
}

func (b *chordBlockParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
}

func (b *chordBlockParser) CanInterruptParagraph() bool {
	return true
}

func (b *chordBlockParser) CanAcceptIndentedLine() bool {
	return true
}

// NewChordBlockParser returns a BlockParser that can parse Tabdown blocks.
func NewChordBlockParser() parser.BlockParser {
	return defaultChordBlockParser
}

type chordParser struct {
}

var defaultChordParser = &chordParser{}

// NewChordParser returns a new  InlineParser that can parse
// Chords in a ChordBlock.
// This parser must take precedence over the parser.LinkParser.
func NewChordParser() parser.InlineParser {
	return defaultChordParser
}

func (s *chordParser) Trigger() []byte {
	return []byte{'['}
}

var taskListRegexp = regexp.MustCompile(`^\[(.*?)\]\s*`)

func (s *chordParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	// Given AST structure must be like
	// - ChordBlock
	//     (current line)
	if parent.Parent() == nil {
		return nil
	}
	if _, ok := parent.(*ast.ChordBlock); !ok {
		return nil
	}
	nr, _ := block.Position()
	if nr > 0 {
		return nil
	}

	line, _ := block.PeekLine()
	m := taskListRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	value := line[m[2]:m[3]]
	indent := block.LineOffset() + m[2] - 1
	block.Advance(m[1])

	return ast.NewChord(indent, value)
}

func (s *chordParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// ChordBlockHTMLRenderer is a renderer.NodeRenderer implementation that
// renders ChordBlock nodes.
type ChordBlockHTMLRenderer struct {
	html.Config
}

// NewChordBlockHTMLRenderer returns a new ChordBlockHTMLRenderer.
func NewChordBlockHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &ChordBlockHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *ChordBlockHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindChordBlock, r.renderChordBlock)
}

func (r *ChordBlockHTMLRenderer) renderChordBlock(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	return gast.WalkContinue, nil
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
			util.Prioritized(NewChordBlockParser(), 0),
		),
		parser.WithInlineParsers(
			util.Prioritized(NewChordParser(), 0),
		),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewChordBlockHTMLRenderer(), 500),
	))
}
