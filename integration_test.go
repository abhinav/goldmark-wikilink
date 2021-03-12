package wikilink_test

import (
	"testing"

	wikilink "github.com/abhinav/goldmark-wikilink"
	"github.com/yuin/goldmark"
	goldtestutil "github.com/yuin/goldmark/testutil"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	goldtestutil.DoTestCaseFile(
		goldmark.New(goldmark.WithExtensions(&wikilink.Extender{})),
		"testdata/integration_test.txt",
		t,
	)
}
