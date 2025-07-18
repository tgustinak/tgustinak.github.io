package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

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
		fmt.Println("Watching for file changes...")
		err := watcher.Watch(*contentDir, func() error {
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

		minifiedHtml, err := minify.Minify([]byte(parsed.HTMLOutput))
		if err != nil {
			return err
		}

		outputFile := filepath.Base(path[:len(path)-3]) + ".html"
		return gen.Generate(map[string]any{
			"Title":       meta.Title,
			"Date":        meta.Date,
			"Tags":        meta.Tags,
			"Content":     template.HTML(minifiedHtml),
			"Description": meta.Description,
		}, outputFile)
	})
}
