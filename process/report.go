package process

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"html/template"
	
	iiifuri "github.com/go-iiif/go-iiif-uri"
	// "github.com/go-iiif/go-iiif/v6/static/templates/html"
	"github.com/jtacoma/uritemplates"
	"github.com/tidwall/gjson"
)

const REPORTNAME_TEMPLATE string = "process_{sha256_origin}.json"

func DeriveReportNameFromURI(ctx context.Context, u iiifuri.URI, uri_template string) (string, error) {

	report_vars := make(map[string]interface{})

	if strings.Contains(uri_template, "{sha256_origin}") {

		origin := u.Origin()

		h := sha256.New()
		h.Write([]byte(origin))

		suffix := fmt.Sprintf("%x", h.Sum(nil))

		report_vars["sha256_origin"] = suffix
	}

	report_t, err := uritemplates.Parse(uri_template)

	if err != nil {
		return "", err
	}

	return report_t.Expand(report_vars)
}

func GenerateProcessReportHTML(ctx context.Context, report_body []byte) ([]byte, error) {

	type Image struct {
		URI    string
		Height int
		Width  int
	}

	var images = make([]*Image, 0)

	uris_rsp := gjson.GetBytes(report_body, "uris")

	for label, u := range uris_rsp.Map() {

		uri := u.String()
		fname := filepath.Base(uri)

		dims_rsp := gjson.GetBytes(report_body, fmt.Sprintf("dimensions.%s", label))
		dims := dims_rsp.Array()

		w := dims[0].Int()
		h := dims[1].Int()

		im := &Image{
			URI:    fname,
			Height: int(h),
			Width:  int(w),
		}

		images = append(images, im)
	}

	sort.Slice(images, func(i, j int) bool {
		return (images[j].Height * images[j].Width) < (images[i].Height * images[i].Width)
	})

	type HTMLVars struct {
		Images []*Image
	}

	var t *template.Template

	/*
	t, err := html.LoadTemplates(ctx)

	if err != nil {
		return nil, err
	}
	*/
	
	t = t.Lookup("process_report")

	if t == nil {
		return nil, fmt.Errorf("Missing 'process_report' template")
	}

	vars := HTMLVars{
		Images: images,
	}

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	err := t.Execute(wr, vars)

	if err != nil {
		return nil, fmt.Errorf("Failed to render 'process_report' template, %w", err)
	}

	wr.Flush()

	return buf.Bytes(), nil
}
