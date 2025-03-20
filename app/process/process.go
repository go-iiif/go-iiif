package process

/*

	./bin/iiif-process \
		-config-source defaults:// \
		-instructions-source defaults:// \
		-verbose \
		-report \
		-config-images-source-uri file:///usr/local \
		-config-derivatives-cache-uri file:///usr/local/test \
		'idsecret:///IMG_9998.jpg?id=9998&secret=abc&secret_o=def&format=jpg&label=x'

*/

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aaronland/gocloud-blob/bucket"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifaws "github.com/go-iiif/go-iiif/v6/aws"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifprocess "github.com/go-iiif/go-iiif/v6/process"
	"gocloud.dev/blob"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	var report_bucket *blob.Bucket

	if opts.ReportSource != "" {

		b, err := bucket.OpenBucket(ctx, opts.ReportSource)

		if err != nil {
			return fmt.Errorf("Failed to open report bucket, %w", err)
		}

		report_bucket = b
		defer report_bucket.Close()
	}

	driver, err := iiifdriver.NewDriverFromConfig(opts.Config)

	if err != nil {
		return fmt.Errorf("Failed to create new driver from config, %w", err)
	}

	pr, err := iiifprocess.NewIIIFProcessor(opts.Config, driver)

	if err != nil {
		return fmt.Errorf("Failed to create new IIIF processor, %w", err)
	}

	/*
		if generate_report_html {
			slog.Info("-generate-report-html flag is true so automatically setting -report=true")
			report = true
		}
	*/

	process_opts := &ProcessOptions{
		Config:         opts.Config,
		Processor:      pr,
		Driver:         driver,
		Instructions:   opts.Instructions,
		Report:         opts.Report,
		ReportTemplate: opts.ReportTemplate,
		ReportBucket:   report_bucket,
		// GenerateReportHTML: generate_report_html,
	}

	switch opts.Mode {

	case "cli":

		to_process := make([]iiifuri.URI, 0)

		for _, str_uri := range opts.Paths {

			u, err := iiifuri.NewURI(ctx, str_uri)

			if err != nil {
				return fmt.Errorf("URI Func for '%s' failed: %w", str_uri, err)
			}

			to_process = append(to_process, u)
		}

		err = ProcessMany(ctx, process_opts, to_process...)

		if err != nil {
			return fmt.Errorf("Failed to process many, %w", err)
		}

	case "fsnotify":

		images_source := opts.Config.Images.Source.Path

		u, err := url.Parse(images_source)

		if err != nil {
			return err
		}

		if u.Scheme != "file" {
			return errors.New("Invalid image source for -mode fsnotify")
		}

		root := u.Path

		logger := slog.Default()
		logger = logger.With("root", root)

		logger.Info("Watching filesystem")

		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			return err
		}

		defer watcher.Close()

		done := make(chan bool)
		wg := new(sync.WaitGroup)

		go func() {

			for {
				select {
				case event, ok := <-watcher.Events:

					if !ok {
						return
					}

					if event.Op == fsnotify.Create {

						abs_path := event.Name

						rel_path := strings.Replace(abs_path, root, "", 1)
						rel_path = strings.TrimLeft(rel_path, "/")

						u, err := iiifuri.NewURI(ctx, rel_path)

						if err != nil {
							logger.Warn("Failed to parse path", "rel path", rel_path, "abs path", abs_path, "error", err)
							continue
						}

						err = ProcessMany(ctx, process_opts, u)

						if err != nil {
							logger.Warn("Failed to process path", "rel path", rel_path, "uri", u, "error", err)
							continue
						}
					}

				case err, ok := <-watcher.Errors:

					if !ok {
						return
					}

					logger.Error("Watch error", "error", err)
				}
			}
		}()

		err = watcher.Add(root)

		if err != nil {
			return err
		}

		<-done

		wg.Wait()

	case "lambda":

		handler := func(ctx context.Context, ev iiifaws.Event) error {

			to_process := make([]iiifuri.URI, 0)

			for _, r := range ev.Records {

				s3_entity := r.S3
				s3_obj := s3_entity.Object
				s3_key := s3_obj.Key

				s3_fname := filepath.Base(s3_key)

				u, err := iiifuri.NewURI(ctx, s3_fname)

				if err != nil {
					return err
				}

				to_process = append(to_process, u)
			}

			err = ProcessMany(ctx, process_opts, to_process...)

			if err != nil {
				return err
			}

			return nil
		}

		aws_lambda.Start(handler)

	default:
		return errors.New("Unsupported mode")
	}

	return nil
}
