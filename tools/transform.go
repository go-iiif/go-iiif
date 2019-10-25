package tools

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/gocloud-blob-bucket"
	aws_events "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiiflevel "github.com/go-iiif/go-iiif/level"
	iiifsource "github.com/go-iiif/go-iiif/source"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"gocloud.dev/blob"
	"io/ioutil"
	"log"
	"path/filepath"
)

type TransformTool struct {
	Tool
}

func NewTransformTool() (Tool, error) {

	t := &TransformTool{}
	return t, nil
}

type TransformOptions struct {
	Config         *iiifconfig.Config
	Driver         iiifdriver.Driver
	Transformation *iiifimage.Transformation
	SourceBucket   *blob.Bucket
	TargetBucket   *blob.Bucket
}

func TransformMany(ctx context.Context, opts *TransformOptions, fnames ...string) error {

	for _, fname := range fnames {

		err := Transform(ctx, opts, fname)

		if err != nil {
			return err
		}
	}

	return nil
}

func Transform(ctx context.Context, opts *TransformOptions, fname string) error {

	fh, err := opts.SourceBucket.NewReader(ctx, fname, nil)

	if err != nil {
		return err
	}

	defer fh.Close()

	if !opts.Transformation.HasTransformation() {
		return errors.New("No transformation")
	}

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return err
	}

	source, err := iiifsource.NewMemorySource(body)

	if err != nil {
		return err
	}

	image, err := opts.Driver.NewImageFromConfigWithSource(opts.Config, source, fname)

	if err != nil {
		return err
	}

	err = image.Transform(opts.Transformation)

	if err != nil {
		return err
	}

	wr, err := opts.TargetBucket.NewWriter(ctx, fname, nil)

	if err != nil {

		return err
	}

	_, err = wr.Write(image.Body())

	if err != nil {
		return err
	}

	err = wr.Close()

	if err != nil {
		return err
	}

	return nil
}

func (t *TransformTool) Run(ctx context.Context) error {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file. DEPRECATED - please use -config_source and -config name.")

	var config_source = flag.String("config-source", "", "A valid Go Cloud bucket URI where your go-iiif config file is located.")
	var config_name = flag.String("config-name", "config.json", "The name of your go-iiif config file.")

	var region = flag.String("region", "full", "A valid IIIF 2.0 region value.")
	var size = flag.String("size", "full", "A valid IIIF 2.0 size value.")
	var rotation = flag.String("rotation", "0", "A valid IIIF 2.0 rotation value.")
	var quality = flag.String("quality", "default", "A valid IIIF 2.0 quality value.")
	var format = flag.String("format", "jpg", "A valid IIIF 2.0 format value.")

	var source_path = flag.String("source", "file:///", "A valid Go Cloud bucket URI where the source file to transform is located.")
	var target_path = flag.String("target", "file:///", "A valid Go Cloud bucket URI where the transformed file should be written.")

	var mode = flag.String("mode", "cli", "Valid modes are: cli, lambda.")

	flag.Parse()

	err := flags.SetFlagsFromEnvVars("IIIF_TRANSFORM")

	if err != nil {
		return err
	}

	if *cfg != "" {

		log.Println("-config flag is deprecated. Please use -config-source and -config-name (setting them now).")

		abs_config, err := filepath.Abs(*cfg)

		if err != nil {
			return err
		}

		*config_name = filepath.Base(abs_config)
		*config_source = fmt.Sprintf("file://%s", filepath.Dir(abs_config))
	}

	config_bucket, err := bucket.OpenBucket(ctx, *config_source)

	if err != nil {
		return err
	}

	config, err := iiifconfig.NewConfigFromBucket(ctx, config_bucket, *config_name)

	if err != nil {
		return err
	}

	driver, err := iiifdriver.NewDriverFromConfig(config)

	if err != nil {
		return err
	}

	// TO DO DEFAULT TO source/target FROM config BUT CHECK FOR OVERRIDE IN *source/target_path ARGS

	source_bucket, err := bucket.OpenBucket(ctx, *source_path)

	if err != nil {
		return err
	}

	target_bucket, err := bucket.OpenBucket(ctx, *target_path)

	if err != nil {
		return err
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "http://127.0.0.1")

	if err != nil {
		return err
	}

	transformation, err := iiifimage.NewTransformation(level, *region, *size, *rotation, *quality, *format)

	if err != nil {
		return err
	}

	transform_opts := &TransformOptions{
		Config:         config,
		Driver:         driver,
		Transformation: transformation,
		SourceBucket:   source_bucket,
		TargetBucket:   target_bucket,
	}

	to_transform := make([]string, 0)

	switch *mode {

	case "cli":

		to_transform = flag.Args()

		err = TransformMany(ctx, transform_opts, to_transform...)

		if err != nil {
			return err
		}

	case "lambda":

		handler := func(ctx context.Context, ev aws_events.S3Event) error {

			for _, r := range ev.Records {

				s3_entity := r.S3
				s3_obj := s3_entity.Object
				s3_key := s3_obj.Key

				s3_fname := filepath.Base(s3_key)
				to_transform = append(to_transform, s3_fname)
			}

			err := TransformMany(ctx, transform_opts, to_transform...)

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
