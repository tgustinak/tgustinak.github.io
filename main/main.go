package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tgustinak.github.io/internal/generator"
	"tgustinak.github.io/internal/minify"
	"tgustinak.github.io/internal/parser"
	"tgustinak.github.io/internal/watcher"
)

type siteConfig struct {
	BaseURL string
	Author  string
	Year    int
}

func (s siteConfig) canonical(outputFile string) string {
	if outputFile == "index.html" {
		return s.BaseURL + "/"
	}

	return s.BaseURL + "/" + filepath.ToSlash(outputFile)
}

func main() {
	contentDir := flag.String("content", "content", "Content directory path")
	templateDir := flag.String("templates", "templates", "Templates directory path")
	outputDir := flag.String("output", ".", "Output directory path")
	baseURL := flag.String("base-url", "https://tgustinak.github.io", "Site base URL for canonical links and sitemap")
	author := flag.String("author", "Tomáš Gustiňák", "Site author name")
	watch := flag.Bool("watch", false, "Watch for file changes")

	flag.Parse()

	site := siteConfig{
		BaseURL: strings.TrimSuffix(*baseURL, "/"),
		Author:  *author,
		Year:    time.Now().Year(),
	}

	gen := generator.NewGenerator(*templateDir, *outputDir)

	if err := build(*contentDir, gen, site); err != nil {
		log.Fatal(err)
	}

	if *watch {
		fmt.Println("Watching for file changes... (Ctrl+C to stop)")
		err := watcher.Watch([]string{*contentDir, *templateDir}, func() error {
			gen.Invalidate()

			return build(*contentDir, gen, site)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func build(contentDir string, gen *generator.Generator, site siteConfig) error {
	pages, err := processFiles(contentDir, gen, site)
	if err != nil {
		return err
	}

	return writeExtras(gen, site, pages)
}

func processFiles(contentDir string, gen *generator.Generator, site siteConfig) ([]string, error) {
	var pages []string

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		meta, content, err := parser.ParseFrontmatter(content)
		if err != nil {
			return err
		}

		parsed := parser.ParseMarkdown(content)

		if meta == nil {
			meta = &parser.Frontmatter{}
		}

		rel, err := filepath.Rel(contentDir, path)
		if err != nil {
			return err
		}

		base := strings.TrimSuffix(filepath.Base(rel), filepath.Ext(rel))
		outputFile := filepath.Join(filepath.Dir(rel), base+".html")

		page, err := gen.Render(pageData(site, outputFile, meta.Title, meta.Description, meta.Tags, meta.Date, template.HTML(parsed.HTMLOutput)))
		if err != nil {
			return err
		}

		minifiedPage, err := minify.Minify([]byte(page))
		if err != nil {
			return err
		}

		if err := gen.Write(outputFile, string(minifiedPage)); err != nil {
			return err
		}

		pages = append(pages, filepath.ToSlash(outputFile))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func pageData(site siteConfig, outputFile string, title string, description string, tags []string, date string, content template.HTML) map[string]any {
	return map[string]any{
		"Title":       title,
		"Date":        date,
		"Tags":        tags,
		"Content":     content,
		"Description": description,
		"Site": map[string]any{
			"BaseURL":   site.BaseURL,
			"Canonical": site.canonical(outputFile),
			"Author":    site.Author,
			"Year":      site.Year,
		},
	}
}

func writeExtras(gen *generator.Generator, site siteConfig, pages []string) error {
	robots := "User-agent: *\nAllow: /\nSitemap: " + site.BaseURL + "/sitemap.xml\n"
	if err := gen.Write("robots.txt", robots); err != nil {
		return err
	}

	var sitemap strings.Builder
	sitemap.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sitemap.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for _, page := range pages {
		loc := site.BaseURL + "/"
		if page != "index.html" {
			loc += page
		}

		sitemap.WriteString("  <url><loc>" + loc + "</loc></url>\n")
	}
	sitemap.WriteString("</urlset>\n")

	if err := gen.Write("sitemap.xml", sitemap.String()); err != nil {
		return err
	}

	notFound, err := gen.Render(pageData(
		site,
		"404.html",
		"Page not found",
		"Page not found",
		nil,
		"",
		template.HTML("<h1>Page not found</h1><p>The page you are looking for does not exist.</p><p><a href=\"index.html\">Back to CV</a></p>"),
	))
	if err != nil {
		return err
	}

	minifiedNotFound, err := minify.Minify([]byte(notFound))
	if err != nil {
		return err
	}

	return gen.Write("404.html", string(minifiedNotFound))
}
