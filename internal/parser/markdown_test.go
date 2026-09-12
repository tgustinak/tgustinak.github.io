package parser

import (
	"strings"
	"testing"
)

func TestParseMarkdownHeadingIDs(t *testing.T) {
	parsed := ParseMarkdown([]byte("# Hello World\n\nSome text.\n"))

	if !strings.Contains(parsed.HTMLOutput, "<h1") {
		t.Fatalf("expected h1 in output, got %q", parsed.HTMLOutput)
	}

	if !strings.Contains(parsed.HTMLOutput, `id="hello-world"`) {
		t.Fatalf("expected auto heading id, got %q", parsed.HTMLOutput)
	}
}
