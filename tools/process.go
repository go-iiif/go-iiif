package tools

// ./bin/iiif-process -config config.json -instructions instructions.json -uri avocado.png
// {"avocado.png":{"b":"avocado.png/full/!2048,1536/0/color.jpg","d":"avocado.png/-1,-1,320,320/full/0/dither.jpg","o":"avocado.png/full/full/-1/color.jpg"}}

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/gocloud-blob-bucket"
	aws_events "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/cache"
	"github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	"github.com/go-iiif/go-iiif/process"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
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
	Config       *config.Config
	Driver       iiifdriver.Driver
	Processor    process.Processor
	Instructions process.IIIFInstructionSet
	Report       bool
	ReportName   string
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
				root := filepath.Dir(target)

				key := filepath.Join(root, opts.ReportName)
				wg.Add(1)

				go func() {

					defer wg.Done()
					err := report_processing(opts.Config, key, rsp)

					if err != nil {
						log.Printf("Unable to write process report %s, %s", key, err)
					}
				}()
			}
		}

		results[origin] = rsp
	}

	wg.Wait()

	return &results, nil
}

func (t *ProcessTool) Run(ctx context.Context) error {

	var iiif_config = flag.String("config", "", "Path to a valid go-iiif config file. DEPRECATED - please use -config_source and -config name.")
	var instructions = flag.String("instructions", "", "Path to a valid go-iiif processing instructions file. DEPRECATED - please use -instructions-source and -instructions-name.")

	var config_source = flag.String("config-source", "", "A valid Go Cloud bucket URI where your go-iiif config file is located.")
	var config_name = flag.String("config-name", "config.json", "The name of your go-iiif config file.")

	var instructions_source = flag.String("instructions-source", "", "A valid Go Cloud bucket URI where your go-iiif instructions file is located.")
	var instructions_name = flag.String("instructions-name", "instructions.json", "The name of your go-iiif instructions file.")

	var report = flag.Bool("report", false, "Store a process report (JSON) for each URI in the cache tree.")
	var report_name = flag.String("report-name", "process.json", "The filename for process reports. Default is 'process.json' as in '${URI}/process.json'.")

	mode := flag.String("mode", "cli", "Valid modes are: cli, fsnotify, lambda.")

	flag.Parse()

	err := flags.SetFlagsFromEnvVars("IIIF_PROCESS")

	if err != nil {
		return err
	}

	if *iiif_config != "" {

		log.Println("-config flag is deprecated. Please use -config-source and -config-name (setting them now).")

		abs_config, err := filepath.Abs(*iiif_config)

		if err != nil {
			return err
		}

		*config_name = filepath.Base(abs_config)
		*config_source = fmt.Sprintf("file://%s", filepath.Dir(abs_config))
	}

	if *instructions != "" {

		log.Println("-instructions flag is deprecated. Please use -instructions-source and -instructions-name (setting them now).")

		abs_instructions, err := filepath.Abs(*instructions)

		if err != nil {
			return err
		}

		*instructions_name = filepath.Base(abs_instructions)
		*instructions_source = fmt.Sprintf("file://%s", filepath.Dir(abs_instructions))
	}

	if *config_source == "" {
		return errors.New("Required -config-source flag is empty.")
	}

	if *instructions_source == "" {
		return errors.New("Required -instructions-source flag is empty.")
	}

	config_bucket, err := bucket.OpenBucket(ctx, *config_source)

	if err != nil {
		return err
	}

	cfg, err := config.NewConfigFromBucket(ctx, config_bucket, *config_name)

	if err != nil {
		return err
	}

	instructions_bucket, err := bucket.OpenBucket(ctx, *instructions_source)

	if err != nil {
		return err
	}

	instructions_set, err := process.ReadInstructionsFromBucket(ctx, instructions_bucket, *instructions_name)

	if err != nil {
		return err
	}

	driver, err := iiifdriver.NewDriverFromConfig(cfg)

	if err != nil {
		return err
	}

	pr, err := process.NewIIIFProcessor(cfg, driver)

	if err != nil {
		return err
	}

	process_opts := &ProcessOptions{
		Config:       cfg,
		Processor:    pr,
		Driver:       driver,
		Instructions: instructions_set,
		Report:       *report,
		ReportName:   *report_name,
	}

	switch *mode {

	case "cli":

		to_process := make([]iiifuri.URI, 0)

		for _, str_uri := range flag.Args() {

			u, err := t.URIFunc(str_uri)

			if err != nil {
				log.Fatal(err)
			}

			to_process = append(to_process, u)
		}

		err = ProcessMany(ctx, process_opts, to_process...)

		if err != nil {
			return err
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

		handler := func(ctx context.Context, ev aws_events.S3Event) error {

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

func report_processing(cfg *config.Config, key string, rsp map[string]interface{}) error {

	dest_cache, err := iiifcache.NewDerivativesCacheFromConfig(cfg)

	if err != nil {
		return err

	}

	enc_rsp, err := json.Marshal(rsp)

	if err != nil {
		return err
	}

	return dest_cache.Set(key, enc_rsp)
}
