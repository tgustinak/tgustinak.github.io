package generator

import (
	"html/template"
	"os"
	"path/filepath"
)

type Generator struct {
	TemplateDir string
	OutputDir   string
}

func NewGenerator(templateDir string, outputDir string) *Generator {
	return &Generator{
		TemplateDir: templateDir,
		OutputDir:   outputDir,
	}
}

func (g *Generator) Generate(data any, outputFile string) error {
	if err := os.MkdirAll(g.OutputDir, 0755); err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(filepath.Join(g.TemplateDir, "default.html"))
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(g.OutputDir, outputFile))
	if err != nil {
		return err
	}
	defer out.Close()

	return tmpl.Execute(out, data)
}
