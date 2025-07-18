package parser

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
}

func ParseFrontmatter(content []byte) (*Frontmatter, []byte, error) {
	parts := bytes.Split(content, []byte("---"))
	if len(parts) < 3 {
		return nil, content, nil
	}

	var meta Frontmatter
	err := yaml.Unmarshal(parts[1], &meta)
	if err != nil {
		return nil, content, err
	}

	return &meta, bytes.Join(parts[2:], []byte("---")), nil
}
