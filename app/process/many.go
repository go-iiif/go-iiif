package process

import (
	"context"
	"log/slog"
	"net/url"
	"path/filepath"
	"sync"

	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifprocess "github.com/go-iiif/go-iiif/v6/process"
	"gocloud.dev/blob"
)

type ProcessResultsReport map[string]interface{}

type ProcessOptions struct {
	Config             *iiifconfig.Config
	Driver             iiifdriver.Driver
	Processor          iiifprocess.Processor
	Instructions       iiifprocess.IIIFInstructionSet
	Report             bool
	ReportTemplate     string
	ReportBucket       *blob.Bucket
	GenerateReportHTML bool
}

func ProcessMany(ctx context.Context, opts *ProcessOptions, uris ...iiifuri.URI) error {

	_, err := ProcessManyWithReport(ctx, opts, uris...)
	return err
}

func ProcessManyWithReport(ctx context.Context, opts *ProcessOptions, uris ...iiifuri.URI) (*ProcessResultsReport, error) {

	results := make(ProcessResultsReport)

	wg := new(sync.WaitGroup)

	for _, uri := range uris {

		origin := uri.Origin()

		logger := slog.Default()
		logger = logger.With("origin", origin)

		rsp, err := iiifprocess.ParallelProcessURIWithInstructionSet(opts.Config, opts.Driver, opts.Processor, opts.Instructions, uri)

		if err != nil {
			return &results, err
		}

		if opts.Report {

			uri_opts := &url.Values{}
			uri_opts.Set("format", "jpg") // this is made up (and not necessarily part of the instructions file)
			uri_opts.Set("label", "x")    // this is made up (and not necessarily part of the instructions file)

			target, err := uri.Target(uri_opts)

			if err != nil {
				logger.Error("Unable to generate target URL for report", "error", err)
			} else {

				report_name, err := iiifprocess.DeriveReportNameFromURI(ctx, uri, opts.ReportTemplate)

				if err == nil {

					var root string

					switch uri.(type) {
					case *iiifuri.IdSecretURI:
						root = filepath.Dir(target)
					default:
						root = target
					}

					key := filepath.Join(root, report_name)
					wg.Add(1)

					go func() {

						defer wg.Done()
						err := report_processing(ctx, opts, key, rsp)

						if err != nil {
							logger.Error("Unable to write process report", "key", key, "error", err)
						}
					}()

				} else {
					logger.Error("Unable to generate report name", "error", err)
				}
			}
		}

		results[origin] = rsp
	}

	wg.Wait()

	return &results, nil
}
