package wikilink

// DefaultResolver is a minimal wiklink resolver that resolves to HTML pages
// relative to the source page.
//
// For example,
//
//  [[Foo]]      // => "Foo.html"
//  [[Foo bar]]  // => "Foo bar.html"
//  [[foo/Bar]]  // => "foo/Bar.html"
var DefaultResolver Resolver = defaultResolver{}

// Resolver resolves pages referenced by wikilinks to their destinations.
type Resolver interface {
	// ResolveWikilink returns the address of the page that the provided
	// wikilink points to. The destination will be URL-escaped before
	// being placed into a link.
	//
	// If ResolveWikilink returns a non-nil error, rendering will be
	// halted.
	ResolveWikilink(*Node) (destination []byte, err error)
}

var _html = []byte(".html")

type defaultResolver struct{}

func (defaultResolver) ResolveWikilink(n *Node) ([]byte, error) {
	dest := make([]byte, len(n.Target)+len(_html))
	copy(dest, n.Target)
	copy(dest[len(n.Target):], _html)
	return dest, nil
}
