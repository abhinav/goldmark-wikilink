package wikilink

import "go.abhg.dev/goldmark/wikilink"

// Kind is the kind of the wikilink AST node.
var Kind = wikilink.Kind

// Node is a Wikilink AST node. Wikilinks have two components: the target and
// the label.
//
// The target is the page to which this link points, and the label is the text
// that displays for this link.
//
// For links in the following form, the label and the target are the same.
//
//	[[Foo bar]]
//
// For links in the following form, the target is the portion of the link to
// the left of the "|", and the label is the portion to the right.
//
//	[[Foo bar|baz qux]]
type Node = wikilink.Node
