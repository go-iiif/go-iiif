package tools

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	"github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	"github.com/go-iiif/go-iiif/v6/process"
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
	Config         *config.Config
	Driver         iiifdriver.Driver
	Processor      process.Processor
	Instructions   process.IIIFInstructionSet
	Report         bool
	ReportTemplate string
	ReportBucket   *blob.Bucket
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

		rsp, err := process.ParallelProcessURIWithInstructionSet(opts.Config, opts.Driver, opts.Processor, opts.Instructions, uri)

		if err != nil {
			return &results, err
		}

		if opts.Report {

			uri_opts := &url.Values{}
			uri_opts.Set("format", "jpg") // this is made up (and not necessarily part of the instructions file)
			uri_opts.Set("label", "x")    // this is made up (and not necessarily part of the instructions file)

			target, err := uri.Target(uri_opts)

			if err != nil {
				log.Printf("Unable to generate target URL for report %s", err)
			} else {

				report_name, err := process.DeriveReportNameFromURI(ctx, uri, opts.ReportTemplate)

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
							log.Printf("Unable to write process report %s, %s", key, err)
						}
					}()

				} else {
					log.Printf("Unable to generate report name %s", err)
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

	err := AppendCommonProcessToolFlags(ctx, fs)

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
	fs.String("report-template", process.REPORTNAME_TEMPLATE, "A valid URI template for generating process report filenames.")
	fs.String("report-source", "", "A valid Go Cloud bucket URI where your report file will be saved. If empty reports will be stored alongside derivative (or cached) images.")

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

	config_source, err := lookup.StringVar(fs, "config-source")

	if err != nil {
		return fmt.Errorf("Failed to lookup -config-source flag, %w", err)
	}

	config_name, err := lookup.StringVar(fs, "config-name")

	if err != nil {
		return fmt.Errorf("Failed to lookup -config-name flag, %w", err)
	}

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

	mode, err := lookup.StringVar(fs, "mode")

	if err != nil {
		return err
	}

	if config_source == "" {
		return errors.New("Required -config-source flag is empty.")
	}

	if instructions_source == "" {
		return errors.New("Required -instructions-source flag is empty.")
	}

	config_bucket, err := blob.OpenBucket(ctx, config_source)

	if err != nil {
		return fmt.Errorf("Failed to open config bucket, %w", err)
	}

	defer config_bucket.Close()

	cfg, err := config.NewConfigFromBucket(ctx, config_bucket, config_name)

	if err != nil {
		return fmt.Errorf("Failed to create new config from bucket %s (%s), %w", config_bucket, config_name, err)
	}

	instructions_bucket, err := blob.OpenBucket(ctx, instructions_source)

	if err != nil {
		return fmt.Errorf("Failed to open instructions bucket, %w", err)
	}

	defer instructions_bucket.Close()

	var report_bucket *blob.Bucket

	if report_source != "" {

		b, err := blob.OpenBucket(ctx, report_source)

		if err != nil {
			return fmt.Errorf("Failed to open report bucket, %w", err)
		}

		report_bucket = b
		defer report_bucket.Close()
	}

	instructions_set, err := process.ReadInstructionsFromBucket(ctx, instructions_bucket, instructions_name)

	if err != nil {
		return fmt.Errorf("Failed to read instructions from bucket, %w", err)
	}

	driver, err := iiifdriver.NewDriverFromConfig(cfg)

	if err != nil {
		return fmt.Errorf("Failed to create new driver from config, %w", err)
	}

	pr, err := process.NewIIIFProcessor(cfg, driver)

	if err != nil {
		return fmt.Errorf("Failed to create new IIIF processor, %w", err)
	}

	process_opts := &ProcessOptions{
		Config:         cfg,
		Processor:      pr,
		Driver:         driver,
		Instructions:   instructions_set,
		Report:         report,
		ReportTemplate: report_template,
		ReportBucket:   report_bucket,
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

		log.Printf("Watching '%s'\n", root)

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
							log.Printf("Failed to parse path '%s' (%s)', %s\n", rel_path, abs_path, err)
							continue
						}

						err = ProcessMany(ctx, process_opts, u)

						if err != nil {
							log.Printf("Failed to process '%s' ('%s'), %s", rel_path, u, err)
							continue
						}
					}

				case err, ok := <-watcher.Errors:

					if !ok {
						return
					}

					log.Println("error:", err)
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

	enc_rsp, err := json.Marshal(rsp)

	if err != nil {
		return err
	}

	if opts.ReportBucket == nil {

		cfg := opts.Config

		dest_cache, err := iiifcache.NewDerivativesCacheFromConfig(cfg)

		if err != nil {
			return err

		}

		return dest_cache.Set(key, enc_rsp)
	}

	fname := filepath.Base(key)

	wr, err := opts.ReportBucket.NewWriter(ctx, fname, nil)

	if err != nil {
		return err
	}

	_, err = wr.Write(enc_rsp)

	if err != nil {
		return err
	}

	return wr.Close()
}
