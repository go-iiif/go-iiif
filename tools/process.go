package tools

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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifprocess "github.com/go-iiif/go-iiif/v6/process"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"gocloud.dev/blob"
)

type ProcessTool struct {
	Tool
	URIFunc URIFunc
}

func NewProcessTool() (Tool, error) {
	uri_func := DefaultURIFunc()
	return NewProcessToolWithURIFunc(uri_func)
}

func NewProcessToolWithURIFunc(uri_func URIFunc) (Tool, error) {

	t := &ProcessTool{
		URIFunc: uri_func,
	}

	return t, nil
}

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

func ProcessToolFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flag.NewFlagSet("process", flag.ExitOnError)

	err := AppendCommonFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	err = AppendCommonProcessToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	err = AppendProcessToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func AppendCommonProcessToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	err := AppendCommonConfigFlags(ctx, fs)

	if err != nil {
		return err
	}

	err = AppendCommonInstructionsFlags(ctx, fs)

	if err != nil {
		return err
	}

	err = AppendCommonToolModeFlags(ctx, fs)

	if err != nil {
		return err
	}

	return nil
}

func AppendProcessToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.Bool("report", false, "Store a process report (JSON) for each URI in the cache tree.")
	fs.String("report-template", iiifprocess.REPORTNAME_TEMPLATE, "A valid URI template for generating process report filenames.")
	fs.String("report-source", "", "A valid Go Cloud bucket URI where your report file will be saved. If empty reports will be stored alongside derivative (or cached) images.")

	fs.Bool("generate-report-html", false, "Generate an HTML page showing all the images listed in a process report.")

	return nil
}

func (t *ProcessTool) Run(ctx context.Context) error {

	fs, err := ProcessToolFlagSet(ctx)

	if err != nil {
		return err
	}

	flagset.Parse(fs)

	err = flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return err
	}

	return t.RunWithFlagSet(ctx, fs)
}

func (t *ProcessTool) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	paths := fs.Args()
	return t.RunWithFlagSetAndPaths(ctx, fs, paths...)
}

func (t *ProcessTool) RunWithFlagSetAndPaths(ctx context.Context, fs *flag.FlagSet, paths ...string) error {

	instructions_source, err := lookup.StringVar(fs, "instructions-source")

	if err != nil {
		return fmt.Errorf("Failed to lookup -instructions-source flag, %w", err)
	}

	instructions_name, err := lookup.StringVar(fs, "instructions-name")

	if err != nil {
		return fmt.Errorf("Failed to lookup -instructions-name flag, %w", err)
	}

	report, err := lookup.BoolVar(fs, "report")

	if err != nil {
		return fmt.Errorf("Failed to lookup -report flag, %w", err)
	}

	report_source, err := lookup.StringVar(fs, "report-source")

	if err != nil {
		return err
	}

	report_template, err := lookup.StringVar(fs, "report-template")

	if err != nil {
		return err
	}

	generate_report_html, err := lookup.BoolVar(fs, "generate-report-html")

	if err != nil {
		return fmt.Errorf("Failed to lookup -generate-report-html flag, %w", err)
	}

	mode, err := lookup.StringVar(fs, "mode")

	if err != nil {
		return err
	}

	verbose, err := lookup.BoolVar(fs, "verbose")

	if err != nil {
		return err
	}

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	cfg, err := iiifconfig.LoadConfigWithFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	instructions_set, err := iiifprocess.LoadInstructions(ctx, instructions_source, instructions_name)

	if err != nil {
		return err
	}

	var report_bucket *blob.Bucket

	if report_source != "" {

		b, err := blob.OpenBucket(ctx, report_source)

		if err != nil {
			return fmt.Errorf("Failed to open report bucket, %w", err)
		}

		report_bucket = b
		defer report_bucket.Close()
	}

	driver, err := iiifdriver.NewDriverFromConfig(cfg)

	if err != nil {
		return fmt.Errorf("Failed to create new driver from config, %w", err)
	}

	pr, err := iiifprocess.NewIIIFProcessor(cfg, driver)

	if err != nil {
		return fmt.Errorf("Failed to create new IIIF processor, %w", err)
	}

	if generate_report_html {
		slog.Info("-generate-report-html flag is true so automatically setting -report=true")
		report = true
	}

	process_opts := &ProcessOptions{
		Config:             cfg,
		Processor:          pr,
		Driver:             driver,
		Instructions:       instructions_set,
		Report:             report,
		ReportTemplate:     report_template,
		ReportBucket:       report_bucket,
		GenerateReportHTML: generate_report_html,
	}

	switch mode {

	case "cli":

		to_process := make([]iiifuri.URI, 0)

		for _, str_uri := range paths {

			u, err := t.URIFunc(str_uri)

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

		images_source := cfg.Images.Source.Path

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

						u, err := t.URIFunc(rel_path)

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

		handler := func(ctx context.Context, ev Event) error {

			to_process := make([]iiifuri.URI, 0)

			for _, r := range ev.Records {

				s3_entity := r.S3
				s3_obj := s3_entity.Object
				s3_key := s3_obj.Key

				s3_fname := filepath.Base(s3_key)

				u, err := t.URIFunc(s3_fname)

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

func report_processing(ctx context.Context, opts *ProcessOptions, key string, rsp map[string]interface{}) error {

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

		return nil
	}

	cfg := opts.Config

	dest_cache, err := iiifcache.NewDerivativesCacheFromConfig(cfg)

	if err != nil {
		return fmt.Errorf("Failed to derive derivatives cache for processing report, %w", err)

	}

	err = dest_cache.Set(key, rsp_body)

	if err != nil {
		return fmt.Errorf("Failed to write report, %w", err)
	}

	slog.Debug("Wrote processing report file", "path", key)

	// START OF HTML version

	if opts.GenerateReportHTML {

		report_html, err := iiifprocess.GenerateProcessReportHTML(ctx, rsp_body)

		html_root := filepath.Dir(key)
		html_path := filepath.Join(html_root, "index.html")

		err = dest_cache.Set(html_path, report_html)

		if err != nil {
			return fmt.Errorf("Failed to write HTML %s, %w", html_path, err)
		}
	}

	return nil
}
