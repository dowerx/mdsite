package doc

import (
	"embed"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed template.html
var htmlTemplateSrc embed.FS
var htmlTemplate *template.Template

var htmlRenderer *html.Renderer

func Init() {
	htmlRenderer = html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})
	htmlTemplate = template.Must(
		template.New("template.html").
			ParseFS(htmlTemplateSrc, "template.html"))
}

type Doc struct {
	template *template.Template
}

func Load(path string) (doc Doc, err error) {
	funcs := template.FuncMap{
		"include": func(path string) (string, error) {
			slog.Info("include", "path", path)
			d, err := Load(path)
			if err != nil {
				return "", err
			}

			return d.Markdown()
		},
		"date": func(layout string) string {
			return time.Now().Format(layout)
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return doc, err
	}

	doc.template, err = template.
		New(filepath.Base(path)).
		Funcs(funcs).
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

func (d *Doc) HTML(head []byte) ([]byte, error) {
	data, err := d.Markdown()
	if err != nil {
		return nil, err
	}

	docRoot := parser.NewWithExtensions(parser.CommonExtensions).Parse([]byte(data))
	body := markdown.Render(docRoot, htmlRenderer)

	var builder strings.Builder
	err = htmlTemplate.Execute(&builder, map[string]string{
		"Head": string(head),
		"Body": string(body),
	})
	return []byte(builder.String()), err
}
