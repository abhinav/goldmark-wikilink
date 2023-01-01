package wikilink

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/ast"
)

func TestRenderer(t *testing.T) {
	t.Parallel()

	t.Run("default resolver", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			desc         string
			give         *Node
			wantEntering string
			wantExiting  string
		}{
			{
				desc: "page",
				give: &Node{
					Target: []byte("foo"),
				},
				wantEntering: `<a href="foo.html">`,
				wantExiting:  `</a>`,
			},
			{
				desc: "image link",
				give: &Node{
					Target: []byte("foo.png"),
				},
				wantEntering: `<a href="foo.png">`,
				wantExiting:  `</a>`,
			},
			{
				desc: "image embed",
				give: &Node{
					Target: []byte("foo.png"),
					Embed:  true,
				},
				wantEntering: `<img src="foo.png">`,
				wantExiting:  ``,
			},
			{
				desc: "image embed url escape",
				give: &Node{
					Target: []byte("my cat picture 1.jpeg"),
					Embed:  true,
				},
				wantEntering: `<img src="my%20cat%20picture%201.jpeg">`,
				wantExiting:  ``,
			},
			{
				desc: "pdf link",
				give: &Node{
					Target: []byte("foo.pdf"),
				},
				wantEntering: `<a href="foo.pdf">`,
				wantExiting:  `</a>`,
			},
			{
				desc: "pdf embed", // unsupported at this time
				give: &Node{
					Target: []byte("foo.pdf"),
					Embed:  true,
				},
				wantEntering: `<a href="foo.pdf">`,
				wantExiting:  `</a>`,
			},
			{
				desc: "page fragment",
				give: &Node{
					Target:   []byte("foo"),
					Fragment: []byte("frag"),
				},
				wantEntering: `<a href="foo.html#frag">`,
				wantExiting:  `</a>`,
			},
			{
				desc: "page fragment embed", // unsupported at this time
				give: &Node{
					Target:   []byte("foo"),
					Fragment: []byte("frag"),
					Embed:    true,
				},
				wantEntering: `<a href="foo.html#frag">`,
				wantExiting:  `</a>`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				var (
					r    Renderer
					buff bytes.Buffer
				)
				w := bufio.NewWriter(&buff)

				_, err := r.Render(w, nil /* source */, tt.give, true /* entering */)
				require.NoError(t, err, "should not fail")
				require.NoError(t, w.Flush(), "flush")

				assert.Equal(t, tt.wantEntering, buff.String(), "output mismatch")
				buff.Reset()

				_, err = r.Render(w, nil /* source */, tt.give, false /* exiting */)
				require.NoError(t, err, "should not fail")
				require.NoError(t, w.Flush(), "flush")

				assert.Equal(t, tt.wantExiting, buff.String(), "output mismatch")
			})
		}
	})

	t.Run("custom resolver", func(t *testing.T) {
		t.Parallel()

		var (
			buff     bytes.Buffer
			w        = bufio.NewWriter(&buff)
			resolved bool
		)
		defer func() {
			assert.True(t, resolved, "custom resolver was never invoked")
		}()

		n := &Node{Target: []byte("foo")}
		r := Renderer{
			Resolver: resolverFunc(func(n *Node) ([]byte, error) {
				assert.False(t, resolved, "resolver invoked too many times")
				resolved = true

				assert.Equal(t, "foo", string(n.Target), "target mismatch")
				return []byte("bar.html"), nil
			}),
		}

		_, err := r.Render(w, nil /* source */, n, true /* entering */)
		require.NoError(t, err, "should not fail")
		require.NoError(t, w.Flush(), "flush")

		assert.Equal(t, `<a href="bar.html">`, buff.String(),
			"output mismatch")
	})

	t.Run("no link", func(t *testing.T) {
		t.Parallel()
		var (
			buff bytes.Buffer
			w    = bufio.NewWriter(&buff)
		)

		n := &Node{Target: []byte("foo")}
		r := Renderer{
			Resolver: resolverFunc(noopResolver),
		}

		_, err := r.Render(w, nil /* source */, n, true /* entering */)
		require.NoError(t, err, "should not fail")

		_, err = r.Render(w, nil /* source */, n, false /* entering */)
		require.NoError(t, err, "should not fail")

		require.NoError(t, w.Flush(), "flush")
		assert.Empty(t, buff.String(), "output should be empty")
	})
}

func TestRenderer_IncorrectNode(t *testing.T) {
	t.Parallel()

	var r Renderer
	_, err := r.Render(bufio.NewWriter(io.Discard), nil /* src */, ast.NewText(), true /* enter */)
	require.Error(t, err, "render with incorrect node must fail")
	assert.Contains(t, err.Error(), "unexpected node")
}

func TestRenderer_ResolveError(t *testing.T) {
	t.Parallel()

	r := Renderer{
		Resolver: resolverFunc(func(*Node) ([]byte, error) {
			return nil, errors.New("great sadness")
		}),
	}
	_, err := r.Render(
		bufio.NewWriter(io.Discard),
		nil, // source
		&Node{Target: []byte("foo")},
		true, // entering
	)
	require.Error(t, err, "render with incorrect node must fail")
	assert.Contains(t, err.Error(), "great sadness")
}

func noopResolver(*Node) ([]byte, error) {
	return nil, nil
}
