package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"

	iiifcache "github.com/go-iiif/go-iiif/v7/cache"
	iiifprocess "github.com/go-iiif/go-iiif/v7/process"
)

func GenerateReports(ctx context.Context, opts *ProcessOptions, key string, rsp map[string]interface{}) error {

	rsp_body, err := json.Marshal(rsp)

	if err != nil {
		return fmt.Errorf("Failed to marshal processing report, %w", err)
	}

	if opts.ReportBucket != nil {

		fname := filepath.Base(key)

		wr, err := opts.ReportBucket.NewWriter(ctx, fname, nil)

		if err != nil {
			return fmt.Errorf("Failed to create new writer for processing report, %w", err)
		}

		_, err = wr.Write(rsp_body)

		if err != nil {
			return fmt.Errorf("Failed to write processing report, %w", err)
		}

		err = wr.Close()

		if err != nil {
			return fmt.Errorf("Failed to close processing report after writing, %w", err)
		}

	} else {

		dest_cache, err := iiifcache.NewCache(ctx, opts.Config.Derivatives.Cache.URI)

		if err != nil {
			return fmt.Errorf("Failed to derive derivatives cache for processing report, %w", err)

		}

		err = dest_cache.Set(key, rsp_body)

		if err != nil {
			return fmt.Errorf("Failed to write report, %w", err)
		}

		slog.Debug("Wrote processing report file", "path", key)
	}

	// START OF HTML version

	if opts.GenerateReportHTML {

		dest_cache, err := iiifcache.NewCache(ctx, opts.Config.Derivatives.Cache.URI)

		if err != nil {
			return fmt.Errorf("Failed to derive derivatives cache for processing report, %w", err)

		}

		report_html, err := iiifprocess.GenerateProcessReportHTML(ctx, rsp_body)

		if err != nil {
			return fmt.Errorf("Failed to generate HTML report, %w", err)
		}

		html_root := filepath.Dir(key)
		html_path := filepath.Join(html_root, "index.html")

		err = dest_cache.Set(html_path, report_html)

		if err != nil {
			return fmt.Errorf("Failed to write HTML %s, %w", html_path, err)
		}
	}

	return nil
}
