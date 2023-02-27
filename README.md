# goldmark-wikilink

[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/goldmark/wikilink.svg)](https://pkg.go.dev/go.abhg.dev/goldmark/wikilink)
[![Go](https://github.com/abhinav/goldmark-wikilink/actions/workflows/go.yml/badge.svg)](https://github.com/abhinav/goldmark-wikilink/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-wikilink/branch/main/graph/badge.svg?token=W98KYF8SPE)](https://codecov.io/gh/abhinav/goldmark-wikilink)

goldmark-wikilink is an extension for the [goldmark] Markdown parser that
supports parsing `[[...]]`-style wiki links
and `![[...]]`-style embedded wiki links.

  [goldmark]: http://github.com/yuin/goldmark

**Demo**:
A web-based demonstration of the extension is available at
<https://abhinav.github.io/goldmark-hashtag/demo/>.

## Installation

```bash
go get go.abhg.dev/goldmark/wikilink@latest
```

## Usage

To use goldmark-wikilink, import the `wikilink` package.

```go
import "go.abhg.dev/goldmark/wikilink"
```

Then include the `wiklink.Extender` in the list of extensions
that you build your [`goldmark.Markdown`] with.

  [`goldmark.Markdown`]: https://pkg.go.dev/github.com/yuin/goldmark#Markdown

```go
goldmark.New(
  goldmark.WithExtensions(
    &wiklink.Extender{},
  ),
  // ...
)
```

## Link resolution

By default, wikilinks will be converted to URLs based on the page name,
unless they already have an extension.

    [[Foo]]     => "Foo.html"
    [[Foo bar]] => "Foo bar.html"
    [[Foo.pdf]] => "Foo.pdf"
    [[Foo.png]] => "Foo.png"

You can change this by supplying a custom [`wikilink.Resolver`]
to your `wikilink.Extender` when you install it.

  [`wikilink.Resolver`]: https://pkg.go.dev/go.abhg.dev/goldmark/wikilink#Resolver

```go
goldmark.New(
  goldmark.WithExtensions(
    // ...
    &wiklink.Extender{
      Resolver: myresolver,
    },
  ),
  // ...
)
```

## Embedding images

Use the embedded link form (`![[...]]`) to add images to a document.

    ![[foo.png]]

Add alt text to images with the `![[...|...]]` form:

    ![[foo.png|alt text]]
