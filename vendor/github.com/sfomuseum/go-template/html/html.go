// Package html provides methods for loading HTML (.html) templates with default functions
package html

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"

	"github.com/sfomuseum/go-template/funcs"
)

// LoadTemplates loads HTML (.html) from 't_fs' with default functions assigned.
func LoadTemplates(ctx context.Context, t_fs ...fs.FS) (*template.Template, error) {

	funcs := TemplatesFuncMap()
	t := template.New("html").Funcs(funcs)

	var err error

	for idx, f := range t_fs {

		t, err = t.ParseFS(f, "*.html")

		if err != nil {
			return nil, fmt.Errorf("Failed to load templates from FS at offset %d, %w", idx, err)
		}
	}

	return t, nil
}

// LoadTemplatesExcluding loads HTML (.html) from 't_fs' with default functions assigned excluding
// templates with (template) names matching 'exclude_list'.
func LoadTemplatesExcluding(ctx context.Context, excluding []string, t_fs ...fs.FS) (*template.Template, error) {

	funcs := TemplatesFuncMap()
	t := template.New("html").Funcs(funcs)

	for _, f := range t_fs {

		err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				return fmt.Errorf("Encountered an error walking %s, %w", path, err)
			}

			r, err := f.Open(path)

			if err != nil {
				return fmt.Errorf("Failed to open %s for reading, %w", path, err)
			}

			defer r.Close()

			body, err := io.ReadAll(r)

			if err != nil {
				return fmt.Errorf("Failed to read %s, %w", path, err)
			}

			this_t, err := t.Parse(string(body))

			if err != nil {
				return fmt.Errorf("Failed to parse %s, %w", path, err)
			}

			t_name := this_t.Name()
			exclude_t := false

			for _, name := range excluding {

				if t_name == name {
					exclude_t = true
					break
				}
			}

			if !exclude_t {
				t, _ = t.Parse(string(body))
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("Failed to parse templates, %w", err)
		}

	}

	return t, nil
}

// TemplatesFuncMap() returns a `template.FuncMap` instance with default functions assigned.
func TemplatesFuncMap() template.FuncMap {

	return template.FuncMap{
		// For example: {{ if (IsAvailable "Account" .) }}
		"IsAvailable":      funcs.IsAvailable,
		"Add":              funcs.Add,
		"JoinPath":         funcs.JoinPath,
		"QRCodeB64":        funcs.QRCodeB64,
		"QRCodeDataURI":    funcs.QRCodeDataURI,
		"IsEven":           funcs.IsEven,
		"IsOdd":            funcs.IsOdd,
		"FormatStringTime": funcs.FormatStringTime,
		"FormatUnixTime":   funcs.FormatUnixTime,
		"GjsonGet":         funcs.GjsonGet,
		"StringHasPrefix":  funcs.StringHasPrefix,
	}
}
