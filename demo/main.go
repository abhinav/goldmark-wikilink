// demo implements a WASM module that can be used to format markdown
// with the goldmark-wikilink extension.
package main

import (
	"bytes"
	"syscall/js"

	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/wikilink"
)

func main() {
	js.Global().Set("renderWikilinks", js.FuncOf(func(this js.Value, args []js.Value) any {
		return renderWikilinks(args[0].String()).Encode()
	}))

	select {}
}

type response struct {
	HTML string
}

func (r *response) Encode() js.Value {
	return js.ValueOf(map[string]any{
		"html": r.HTML,
	})
}

func renderWikilinks(markdown string) *response {
	md := goldmark.New(
		goldmark.WithExtensions(
			&wikilink.Extender{},
		),
	)

	var buff bytes.Buffer
	md.Convert([]byte(markdown), &buff)
	return &response{
		HTML: buff.String(),
	}
}
