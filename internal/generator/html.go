package generator

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
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

func (g *Generator) Invalidate() {
	g.tmpl = nil
	g.tmplErr = nil
	g.tmplParsed = false
}

func (g *Generator) Render(data any) (string, error) {
	tmpl, err := g.template()
	if err != nil {
		return "", err
	}

	var out strings.Builder
	if err := tmpl.Execute(&out, data); err != nil {
		return "", err
	}

	return out.String(), nil
}

func (g *Generator) Write(outputFile string, content string) error {
	outPath := filepath.Join(g.OutputDir, outputFile)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(outPath, []byte(content), 0o644)
}
