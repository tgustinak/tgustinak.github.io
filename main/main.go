package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"tgustinak.github.io/internal/generator"
	"tgustinak.github.io/internal/minify"
	"tgustinak.github.io/internal/parser"
	"tgustinak.github.io/internal/watcher"
)

func main() {
	contentDir := flag.String("content", "content", "Content directory path")
	templateDir := flag.String("templates", "templates", "Templates directory path")
	outputDir := flag.String("output", ".", "Output directory path")
	watch := flag.Bool("watch", false, "Watch for file changes")

	flag.Parse()

	gen := generator.NewGenerator(*templateDir, *outputDir)

	err := processFiles(*contentDir, gen)
	if err != nil {
		log.Fatal(err)
	}

	if *watch {
		fmt.Println("Watching for file changes... (Ctrl+C to stop)")
		err := watcher.Watch([]string{*contentDir, *templateDir}, func() error {
			gen.Invalidate()

			return processFiles(*contentDir, gen)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func processFiles(contentDir string, gen *generator.Generator) error {
	return filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
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

		page, err := gen.Render(map[string]any{
			"Title":       meta.Title,
			"Date":        meta.Date,
			"Tags":        meta.Tags,
			"Content":     template.HTML(parsed.HTMLOutput),
			"Description": meta.Description,
		})
		if err != nil {
			return err
		}

		minifiedPage, err := minify.Minify([]byte(page))
		if err != nil {
			return err
		}

		return gen.Write(outputFile, string(minifiedPage))
	})
}
