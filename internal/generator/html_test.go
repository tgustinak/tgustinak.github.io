package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testTemplate = `<title>{{.Title}}</title><div>{{.Content}}</div>`

func writeTestTemplate(t *testing.T, dir string, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, "default.html"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRenderIncludesData(t *testing.T) {
	dir := t.TempDir()
	writeTestTemplate(t, dir, testTemplate)

	gen := NewGenerator(dir, t.TempDir())

	out, err := gen.Render(map[string]any{"Title": "Hi", "Content": "body"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "<title>Hi</title>") {
		t.Fatalf("expected title in output, got %q", out)
	}
}

func TestRenderMissingTemplate(t *testing.T) {
	gen := NewGenerator(t.TempDir(), t.TempDir())

	_, err := gen.Render(map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing template, got nil")
	}
}

func TestWriteCreatesSubdirs(t *testing.T) {
	dir := t.TempDir()
	writeTestTemplate(t, dir, testTemplate)

	outDir := t.TempDir()
	gen := NewGenerator(dir, outDir)

	if err := gen.Write("sub/page.html", "hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outDir, "sub", "page.html"))
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "hello" {
		t.Fatalf("expected hello, got %q", content)
	}
}

func TestInvalidateReloadsTemplate(t *testing.T) {
	dir := t.TempDir()
	writeTestTemplate(t, dir, "v1 {{.Title}}")

	gen := NewGenerator(dir, t.TempDir())

	first, err := gen.Render(map[string]any{"Title": "x"})
	if err != nil {
		t.Fatal(err)
	}

	if first != "v1 x" {
		t.Fatalf("expected v1, got %q", first)
	}

	writeTestTemplate(t, dir, "v2 {{.Title}}")

	stale, err := gen.Render(map[string]any{"Title": "x"})
	if err != nil {
		t.Fatal(err)
	}

	if stale != "v1 x" {
		t.Fatalf("expected cached v1, got %q", stale)
	}

	gen.Invalidate()

	fresh, err := gen.Render(map[string]any{"Title": "x"})
	if err != nil {
		t.Fatal(err)
	}

	if fresh != "v2 x" {
		t.Fatalf("expected reloaded v2, got %q", fresh)
	}
}
