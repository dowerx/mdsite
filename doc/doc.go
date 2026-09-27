package doc

import (
	"os"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Doc struct {
	template *template.Template
}

func Load(path string) (doc Doc, err error) {
	data, err := os.ReadFile(path)
	doc.template, err = template.
		New("").
		Funcs(template.FuncMap{
			"include": func(path string) (string, error) {
				d, err := Load(path)
				if err != nil {
					return "", err
				}

				return d.Markdown()
			}}).
		Parse(string(data))

	return doc, err
}

func (d *Doc) Markdown() (string, error) {
	var builder strings.Builder
	err := d.template.Execute(&builder, nil)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (d *Doc) HTML() ([]byte, error) {
	data, err := d.Markdown()
	if err != nil {
		return nil, err
	}

	docRoot := parser.
		NewWithExtensions(parser.CommonExtensions).
		Parse([]byte(data))

	return markdown.Render(docRoot,
		html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})), nil
}
