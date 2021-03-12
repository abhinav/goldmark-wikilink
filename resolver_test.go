package wikilink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type resolverFunc func(*Node) ([]byte, error)

func (f resolverFunc) ResolveWikilink(n *Node) ([]byte, error) {
	return f(n)
}

func TestDefaultResolver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give string
		want string
	}{
		{"foo", "foo.html"},
		{"foo bar", "foo bar.html"},
		{"foo/bar", "foo/bar.html"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.give, func(t *testing.T) {
			t.Parallel()

			got, err := DefaultResolver.ResolveWikilink(&Node{
				Target: []byte(tt.give),
			})
			require.NoError(t, err, "resolve failed")
			assert.Equal(t, tt.want, string(got), "result mismatch")
		})
	}
}
