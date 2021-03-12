package wikilink

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Renderer renders wikilinks as HTML.
//
// Install it on your goldmark Markdown object with Extender, or directly on a
// goldmark Renderer by using the WithNodeRenderers option.
//
//   wikilinkRenderer := util.Prioritized(&wikilink.Renderer{...}, 199)
//   goldmarkRenderer.AddOptions(renderer.WithNodeRenderers(wikilinkRenderer))
type Renderer struct {
	// Resolver determines how destinations for wikilink pages.
	//
	// Uses DefaultResolver if unspecified.
	Resolver Resolver
}

// RegisterFuncs registers wikilink rendering functions with the provided
// goldmark registerer. This teaches goldmark to call us when it encounters a
// wikilink in the AST.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.Render)
}

// Render renders the provided Node. It must be a Wikilink node.
//
// goldmark will call this method if this renderer was registered with it
// using the WithNodeRenderers option.
func (r *Renderer) Render(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, ok := node.(*Node)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *goldmarkwikilink.Node", node)
	}

	if entering {
		dest, err := r.resolve(n)
		if err != nil {
			return ast.WalkStop, fmt.Errorf("resolve %q: %w", n.Target, err)
		}

		w.WriteString(`<a href="`)
		w.Write(util.URLEscape(dest, true /* resolve references */))
		w.WriteString(`">`)
	} else {
		w.WriteString("</a>")
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) resolve(n *Node) ([]byte, error) {
	res := r.Resolver
	if res == nil {
		res = DefaultResolver
	}
	return res.ResolveWikilink(n)
}
