package generator

import (
	"html/template"
	"os"
	"path/filepath"
)

type Generator struct {
	TemplateDir string
	OutputDir   string

	tmpl         *template.Template
	tmplErr      error
	tmplParsed   bool
	templateFile string
}

func NewGenerator(templateDir string, outputDir string) *Generator {
	return &Generator{
		TemplateDir:  templateDir,
		OutputDir:    outputDir,
		templateFile: filepath.Join(templateDir, "default.html"),
	}
}

func (g *Generator) template() (*template.Template, error) {
	if g.tmplParsed {
		return g.tmpl, g.tmplErr
	}

	g.tmpl, g.tmplErr = template.ParseFiles(g.templateFile)
	g.tmplParsed = true

	return g.tmpl, g.tmplErr
}

func (g *Generator) Generate(data any, outputFile string) error {
	outPath := filepath.Join(g.OutputDir, outputFile)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	tmpl, err := g.template()
	if err != nil {
		return err
	}

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return tmpl.Execute(out, data)
}
