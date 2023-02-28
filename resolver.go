package wikilink

import "path/filepath"

// DefaultResolver is a minimal wiklink resolver that resolves wikilinks
// relative to the source page.
//
// It adds ".html" to the end of the target
// if the target does not have an extension.
//
// For example,
//
//	[[Foo]]      // => "Foo.html"
//	[[Foo bar]]  // => "Foo bar.html"
//	[[foo/Bar]]  // => "foo/Bar.html"
//	[[foo.pdf]]  // => "foo.pdf"
//	[[foo.png]]  // => "foo.png"
var DefaultResolver Resolver = defaultResolver{}

// Resolver resolves pages referenced by wikilinks to their destinations.
type Resolver interface {
	// ResolveWikilink returns the address of the page that the provided
	// wikilink points to. The destination will be URL-escaped before
	// being placed into a link.
	//
	// If ResolveWikilink returns a non-nil error, rendering will be
	// halted.
	//
	// If ResolveWikilink returns a nil destination and error, the
	// Renderer will omit the link and render its contents as a regular
	// string.
	ResolveWikilink(*Node) (destination []byte, err error)
}

var _html = []byte(".html")

type defaultResolver struct{}

func (defaultResolver) ResolveWikilink(n *Node) ([]byte, error) {
	dest := make([]byte, len(n.Target)+len(_html)+len(_hash)+len(n.Fragment))
	var i int
	if len(n.Target) > 0 {
		i += copy(dest, n.Target)
		if filepath.Ext(string(n.Target)) == "" {
			i += copy(dest[i:], _html)
		}
	}
	if len(n.Fragment) > 0 {
		i += copy(dest[i:], _hash)
		i += copy(dest[i:], n.Fragment)
	}
	return dest[:i], nil
}
