package parser

import (
	"strings"
	"testing"
)

func TestParseFrontmatterkeepsHorizontalRulesInBody(t *testing.T) {
	in := []byte("---\ntitle: CV\n---\n\n# Hi\n\n---\n\nbody after rule\n")

	meta, body, err := ParseFrontmatter(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Title != "CV" {
		t.Fatalf("expected title CV, got %q", meta.Title)
	}

	if !strings.Contains(string(body), "---") {
		t.Fatalf("expected body to keep the horizontal rule, got %q", body)
	}
}

func TestParseFrontmatterWithoutFrontmatter(t *testing.T) {
	in := []byte("# No FM\n\n---\n\ntext\n")

	meta, body, err := ParseFrontmatter(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta == nil {
		t.Fatal("expected non-nil meta")
	}

	if string(body) != string(in) {
		t.Fatalf("expected body passthrough, got %q", body)
	}
}

func TestParseFrontmatterUnclosed(t *testing.T) {
	in := []byte("---\ntitle: oops\n\n# body")

	_, body, err := ParseFrontmatter(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(body) != string(in) {
		t.Fatalf("expected body passthrough, got %q", body)
	}
}

func TestParseFrontmatterInvalidYAML(t *testing.T) {
	in := []byte("---\ntitle: [unclosed\n---\n\nbody\n")

	_, _, err := ParseFrontmatter(in)
	if err == nil {
		t.Fatal("expected YAML error, got nil")
	}
}

func TestParseFrontmatterCRLF(t *testing.T) {
	in := []byte("---\r\ntitle: CV\r\n---\r\n\r\n# Hi\r\n")

	meta, body, err := ParseFrontmatter(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Title != "CV" {
		t.Fatalf("expected title CV, got %q", meta.Title)
	}

	if !strings.Contains(string(body), "# Hi") {
		t.Fatalf("expected body to contain heading, got %q", body)
	}
}
