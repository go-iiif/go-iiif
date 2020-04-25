package process

import (
	"context"
	"crypto/sha256"
	"fmt"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	"github.com/jtacoma/uritemplates"
	"strings"
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
