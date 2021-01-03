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

var chordRegexp = regexp.MustCompile(`\[(.*?)\]`)

func (b *chordBlockParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if !isChordLine(line) {
		return nil, parser.NoChildren
	}

	chordBlock := ast.NewChordBlock()

	line, _ = reader.PeekLine()
	m := chordRegexp.FindSubmatchIndex(line)
	for m != nil {
		value := line[m[2]:m[3]]
		indent := reader.LineOffset() + m[2] - 1
		chord := ast.NewChord(indent, value)
		chordBlock.AppendChild(chordBlock, chord)
		reader.Advance(m[1] - 1)
		line, _ = reader.PeekLine()
		m = chordRegexp.FindSubmatchIndex(line)
	}
	return chordBlock, parser.NoChildren
}

func (b *chordBlockParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	_, lyrics := reader.PeekLine()
	if node.FirstChild() == nil {
		return parser.Close
	}
	chord := node.FirstChild().(*ast.Chord)
	if chord.Indent > 0 {
		prefix := text.NewSegment(lyrics.Start, lyrics.Start+chord.Indent)
		text := gast.NewTextBlock()
		text.Lines().Append(prefix)
		node.InsertBefore(node, chord, text)
	}
	for {
		start := lyrics.Start + chord.Indent
		if start > lyrics.Stop {
			start = lyrics.Stop
		}
		if chord.NextSibling() == nil {
			segment := text.NewSegment(start, lyrics.Stop)
			chord.Lines().Append(segment)
			break
		}
		nextChord := chord.NextSibling().(*ast.Chord)
		stop := lyrics.Start + nextChord.Indent
		if stop > lyrics.Stop {
			stop = lyrics.Stop
		}
		segment := text.NewSegment(start, stop)
		chord.Lines().Append(segment)
		chord = nextChord
	}

	reader.Advance(lyrics.Len() - 1)
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
	reg.Register(ast.KindChord, r.renderChord)
}

func (r *ChordBlockHTMLRenderer) renderChordBlock(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<div class='chord-block'>")
	} else {
		_, _ = w.WriteString("</div>")
	}
	return gast.WalkContinue, nil
}

func (r *ChordBlockHTMLRenderer) renderChord(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.Chord)
	if entering {
		_, _ = w.WriteString("<span class='chord'><span class='chord-name'>")
		_, _ = w.Write(n.Name)
		_, _ = w.WriteString("</span>")

	} else {
		_, _ = w.WriteString("</span>")
	}
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
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewChordBlockHTMLRenderer(), 500),
	))
}
