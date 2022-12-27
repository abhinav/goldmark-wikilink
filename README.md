[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/goldmark/wikilink.svg)](https://pkg.go.dev/go.abhg.dev/goldmark/wikilink)
[![Go](https://github.com/abhinav/goldmark-wikilink/actions/workflows/go.yml/badge.svg)](https://github.com/abhinav/goldmark-wikilink/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-wikilink/branch/main/graph/badge.svg?token=W98KYF8SPE)](https://codecov.io/gh/abhinav/goldmark-wikilink)

goldmark-wikilink is an extension for the [goldmark] Markdown parser that
supports parsing `[[...]]`-style wiki links.

  [goldmark]: http://github.com/yuin/goldmark

# Usage

To use goldmark-wikilink, import the `wikilink` package.

```go
import "go.abhg.dev/goldmark/wikilink"
```

Then include the `wiklink.Extender` in the list of extensions you build your
[`goldmark.Markdown`] with.

  [`goldmark.Markdown`]: https://pkg.go.dev/github.com/yuin/goldmark#Markdown

```go
goldmark.New(
  goldmark.WithExtensions(
    &wiklink.Extender{},
  ),
  // ...
)
```

goldmark-wikilink provides control over destinations of wikilinks with the
`Resolver` type. Specify a custom `Resolver` to the `Extender` when installing
it.

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

