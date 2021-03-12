package wikilink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		give string

		wantTarget string
		wantLabel  string

		remainder string // unconsumed portion of tt.give
	}{
		{
			desc:       "simple",
			give:       "[[foo]] bar",
			wantTarget: "foo",
			wantLabel:  "foo",
			remainder:  " bar",
		},
		{
			desc:       "spaces",
			give:       "[[foo bar]]baz",
			wantTarget: "foo bar",
			wantLabel:  "foo bar",
			remainder:  "baz",
		},
		{
			desc:       "label",
			give:       "[[foo|bar]]",
			wantTarget: "foo",
			wantLabel:  "bar",
		},
		{
			desc:       "label with spaces",
			give:       "[[foo bar|baz qux]] quux",
			wantTarget: "foo bar",
			wantLabel:  "baz qux",
			remainder:  " quux",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			r := text.NewReader([]byte(tt.give))

			var p Parser
			got := p.Parse(nil /* parent */, r, parser.NewContext())
			require.NotNil(t, got, "expected Node, got nil")

			if n, ok := got.(*Node); assert.True(t, ok, "expected Node, got %T", got) {
				assert.Equal(t, tt.wantTarget, string(n.Target), "target mismatch")
			}

			if assert.Equal(t, 1, got.ChildCount(), "children mismatch") {
				child := got.FirstChild()
				if label, ok := child.(*ast.Text); assert.True(t, ok, "expected Text, got %T", child) {
					assert.Equal(t, tt.wantLabel, string(r.Value(label.Segment)), "label mismatch")
				}
			}

			_, pos := r.Position()
			assert.Equal(t, tt.remainder, string(r.Value(pos)),
				"remaining text does not match")
		})
	}
}
