package wikilink

import "go.abhg.dev/goldmark/wikilink"

// Renderer renders wikilinks as HTML.
//
// Install it on your goldmark Markdown object with Extender, or directly on a
// goldmark Renderer by using the WithNodeRenderers option.
//
//	wikilinkRenderer := util.Prioritized(&wikilink.Renderer{...}, 199)
//	goldmarkRenderer.AddOptions(renderer.WithNodeRenderers(wikilinkRenderer))
type Renderer = wikilink.Renderer
