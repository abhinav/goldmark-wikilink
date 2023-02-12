package wikilink

import "go.abhg.dev/goldmark/wikilink"

// DefaultResolver is a minimal wiklink resolver that resolves to HTML pages
// relative to the source page.
//
// For example,
//
//	[[Foo]]      // => "Foo.html"
//	[[Foo bar]]  // => "Foo bar.html"
//	[[foo/Bar]]  // => "foo/Bar.html"
var DefaultResolver Resolver = wikilink.DefaultResolver

// Resolver resolves pages referenced by wikilinks to their destinations.
type Resolver = wikilink.Resolver
