package wikilink

import "go.abhg.dev/goldmark/wikilink"

// Parser parses wikilinks.
//
// Install it on your goldmark Markdown object with Extender, or install it
// directly on your goldmark Parser by using the WithInlineParsers option.
//
//	wikilinkParser := util.Prioritized(&wikilink.Parser{...}, 199)
//	goldmarkParser.AddOptions(parser.WithInlineParsers(wikilinkParser))
//
// Note that the priority for the wikilink parser must 199 or lower to take
// precednce over the plain Markdown link parser which has a priority of 200.
type Parser = wikilink.Parser
