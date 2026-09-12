package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"tgustinak.github.io/internal/generator"
)

func TestBuildGeneratesSite(t *testing.T) {
	contentDir := t.TempDir()
	templateDir := t.TempDir()
	outputDir := t.TempDir()

	md := "---\ntitle: Test\n---\n\n# Hello\n"
	if err := os.WriteFile(filepath.Join(contentDir, "index.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(contentDir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(contentDir, "sub", "page.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}

	tmpl := `<title>{{.Title}}</title><link rel="canonical" href="{{.Site.Canonical}}"><div>{{.Content}}</div>`
	if err := os.WriteFile(filepath.Join(templateDir, "default.html"), []byte(tmpl), 0o644); err != nil {
		t.Fatal(err)
	}

	gen := generator.NewGenerator(templateDir, outputDir)
	site := siteConfig{BaseURL: "https://example.com", Author: "Tester", Year: 2026}

	if err := build(contentDir, gen, site); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{"index.html", "sub/page.html", "robots.txt", "sitemap.xml", "404.html"} {
		if _, err := os.Stat(filepath.Join(outputDir, want)); err != nil {
			t.Fatalf("expected %s to exist: %v", want, err)
		}
	}

	index, err := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(index), `href=https://example.com/`) {
		t.Fatalf("expected root canonical URL, got %q", index)
	}

	sub, err := os.ReadFile(filepath.Join(outputDir, "sub", "page.html"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(sub), `href=https://example.com/sub/page.html`) {
		t.Fatalf("expected sub canonical URL, got %q", sub)
	}

	sitemap, err := os.ReadFile(filepath.Join(outputDir, "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(sitemap), "<loc>https://example.com/</loc>") {
		t.Fatalf("expected sitemap root entry, got %q", sitemap)
	}

	if !strings.Contains(string(sitemap), "<loc>https://example.com/sub/page.html</loc>") {
		t.Fatalf("expected sitemap sub entry, got %q", sitemap)
	}
}
