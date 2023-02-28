package wikilink

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Parser parses wikilinks.
//
// Install it on your goldmark Markdown object with Extender, or install it
// directly on your goldmark Parser by using the WithInlineParsers option.
//
//	wikilinkParser := util.Prioritized(&wikilink.Parser{...}, 199)
//	goldmarkParser.AddOptions(parser.WithInlineParsers(wikilinkParser))
//
// Note that the priority for the wikilink parser must 199 or lower to take
// precedence over the plain Markdown link parser which has a priority of 200.
type Parser struct{}

var _ parser.InlineParser = (*Parser)(nil)

var (
	_open      = []byte("[[")
	_embedOpen = []byte("![[")
	_pipe      = []byte{'|'}
	_hash      = []byte{'#'}
	_close     = []byte("]]")
)

// Trigger returns characters that trigger this parser.
func (p *Parser) Trigger() []byte {
	return []byte{'!', '['}
}

// Parse parses a wikilink in one of the following forms:
//
//	[[...]]    (simple)
//	![[...]]   (embedded)
//
// Both, simple and embedded wikilinks support the following syntax:
//
//	[[target]]
//	[[target|label]]
//
// If the label is omitted, the target is used as the label.
//
// The target may optionally contain a fragment identifier:
//
//	[[target#fragment]]
func (p *Parser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()
	stop := bytes.Index(line, _close)
	if stop < 0 {
		return nil // must close on the same line
	}

	var embed bool

	switch {
	case bytes.HasPrefix(line, _open):
		seg = text.NewSegment(seg.Start+len(_open), seg.Start+stop)
	case bytes.HasPrefix(line, _embedOpen):
		embed = true
		seg = text.NewSegment(seg.Start+len(_embedOpen), seg.Start+stop)
	default:
		return nil
	}

	n := &Node{Target: block.Value(seg), Embed: embed}
	if idx := bytes.Index(n.Target, _pipe); idx >= 0 {
		n.Target = n.Target[:idx]                // [[ ... |
		seg = seg.WithStart(seg.Start + idx + 1) // | ... ]]
	}

	if len(n.Target) == 0 || seg.Len() == 0 {
		return nil // target and label must not be empty
	}

	// Target may be Foo#Bar, so break them apart.
	if idx := bytes.LastIndex(n.Target, _hash); idx >= 0 {
		n.Fragment = n.Target[idx+1:] // Foo#Bar => Bar
		n.Target = n.Target[:idx]     // Foo#Bar => Foo
	}

	n.AppendChild(n, ast.NewTextSegment(seg))
	block.Advance(stop + 2)
	return n
}
