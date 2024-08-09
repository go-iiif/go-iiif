package html

import (
	"context"
	"embed"
	"html/template"

	sfom_html "github.com/sfomuseum/go-template/html"
)

//go:embed *.html
var FS embed.FS

// LoadTemplates loads the templates in the `html` package's embedded filesystem
// and returns a new `template.Template` instance with support for the (template)
// functions defined in `TemplatesFuncMap`.
func LoadTemplates(ctx context.Context) (*template.Template, error) {
	return sfom_html.LoadTemplates(ctx, FS)
}
