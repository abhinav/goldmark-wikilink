package wikilink

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Renderer renders wikilinks as HTML.
//
// Install it on your goldmark Markdown object with Extender, or directly on a
// goldmark Renderer by using the WithNodeRenderers option.
//
//	wikilinkRenderer := util.Prioritized(&wikilink.Renderer{...}, 199)
//	goldmarkRenderer.AddOptions(renderer.WithNodeRenderers(wikilinkRenderer))
type Renderer struct {
	// Resolver determines destinations for wikilink pages.
	//
	// If a Resolver returns an empty destination, the Renderer will skip
	// the link and render just its contents. That is, instead of,
	//
	//   <a href="foo">bar</a>
	//
	// The renderer will render just the following.
	//
	//   bar
	//
	// Defaults to DefaultResolver if unspecified.
	Resolver Resolver

	once sync.Once // guards init

	// hasDest records whether a node had a destination when we resolved
	// it. This is needed to decide whether a closing </a> must be added
	// when exiting a Node render.
	hasDest sync.Map // *Node => struct{}
}

func (r *Renderer) init() {
	r.once.Do(func() {
		if r.Resolver == nil {
			r.Resolver = DefaultResolver
		}
	})
}

// RegisterFuncs registers wikilink rendering functions with the provided
// goldmark registerer. This teaches goldmark to call us when it encounters a
// wikilink in the AST.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.Render)
}

// Render renders the provided Node. It must be a Wikilink [Node].
//
// goldmark will call this method if this renderer was registered with it
// using the WithNodeRenderers option.
//
// All nodes will be rendered as links (with <a> tags),
// except for embed links (![[..]]) that refer to images.
// Those will be rendered as images (with <img> tags).
func (r *Renderer) Render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	r.init()

	n, ok := node.(*Node)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *wikilink.Node", node)
	}

	if entering {
		return r.enter(w, n, src)
	}

	r.exit(w, n)
	return ast.WalkContinue, nil
}

func (r *Renderer) enter(w util.BufWriter, n *Node, src []byte) (ast.WalkStatus, error) {
	dest, err := r.Resolver.ResolveWikilink(n)
	if err != nil {
		return ast.WalkStop, fmt.Errorf("resolve %q: %w", n.Target, err)
	}
	if len(dest) == 0 {
		return ast.WalkContinue, nil
	}

	img := resolveAsImage(n)
	if !img {
		r.hasDest.Store(n, struct{}{})
		_, _ = w.WriteString(`<a href="`)
		_, _ = w.Write(util.URLEscape(dest, true /* resolve references */))
		_, _ = w.WriteString(`">`)
		return ast.WalkContinue, nil
	}

	_, _ = w.WriteString(`<img src="`)
	_, _ = w.Write(util.URLEscape(dest, true /* resolve references */))
	// The label portion of the link becomes the alt text
	// only if it isn't the same as the target.
	// This way, [[foo.jpg]] does not become alt="foo.jpg",
	// but [[foo.jpg|bar]] does become alt="bar".
	if n.ChildCount() == 1 {
		label := nodeText(src, n.FirstChild())
		if !bytes.Equal(label, n.Target) {
			_, _ = w.WriteString(`" alt="`)
			_, _ = w.Write(util.EscapeHTML(label))
		}
	}
	_, _ = w.WriteString(`">`)
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) exit(w util.BufWriter, n *Node) {
	if _, ok := r.hasDest.LoadAndDelete(n); ok {
		_, _ = w.WriteString("</a>")
	}
}

// returns true if the wikilink should be resolved to an image node
func resolveAsImage(n *Node) bool {
	if !n.Embed {
		return false
	}

	filename := string(n.Target)
	switch ext := filepath.Ext(filename); ext {
	// Common image file types taken from
	// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Image_types
	case ".apng", ".avif", ".gif", ".jpg", ".jpeg", ".jfif", ".pjpeg", ".pjp", ".png", ".svg", ".webp":
		return true
	default:
		return false
	}
}

func nodeText(src []byte, n ast.Node) []byte {
	var buf bytes.Buffer
	writeNodeText(src, &buf, n)
	return buf.Bytes()
}

func writeNodeText(src []byte, dst io.Writer, n ast.Node) {
	switch n := n.(type) {
	case *ast.Text:
		_, _ = dst.Write(n.Segment.Value(src))
	case *ast.String:
		_, _ = dst.Write(n.Value)
	default:
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			writeNodeText(src, dst, c)
		}
	}
}
