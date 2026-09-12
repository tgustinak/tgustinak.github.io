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
	// Only treat leading ---\n...\n--- as frontmatter.
	// A naive bytes.Split on "---" breaks as soon as the markdown
	// body itself contains a horizontal rule.
	trimmed := bytes.TrimLeft(content, "\xef\xbb\xbf \t\r\n")
	if !bytes.HasPrefix(trimmed, []byte("---")) {
		return &Frontmatter{}, content, nil
	}

	rest := trimmed[len("---"):]
	// Skip the newline after opening delimiter; require it.
	if len(rest) == 0 || (rest[0] != '\n' && rest[0] != '\r') {
		return &Frontmatter{}, content, nil
	}

	// Find closing delimiter on its own line: \n---(\n or \r or end).
	closing := -1
	closingLen := 0
	for i := range rest {
		if rest[i] != '\n' {
			continue
		}

		j := i + 1
		if j+3 <= len(rest) && string(rest[j:j+3]) == "---" {
			k := j + 3
			if k >= len(rest) || rest[k] == '\n' || rest[k] == '\r' {
				closing = i + 1
				if k < len(rest) && rest[k] == '\r' {
					closingLen = 4 // ---\r
				} else {
					closingLen = 3 // ---
				}
				break
			}
		}
	}

	if closing == -1 {
		return &Frontmatter{}, content, nil
	}

	var meta Frontmatter
	if err := yaml.Unmarshal(rest[:closing-1], &meta); err != nil {
		return nil, content, err
	}

	body := rest[closing+closingLen:]
	body = bytes.TrimLeft(body, "\r\n")

	return &meta, body, nil
}
