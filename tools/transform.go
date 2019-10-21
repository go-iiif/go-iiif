package tools

import (
	"context"
	"errors"
	"flag"
	"github.com/aaronland/gocloud-blob-bucket"
	aws_events "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiiflevel "github.com/go-iiif/go-iiif/level"
	iiifsource "github.com/go-iiif/go-iiif/source"
	"gocloud.dev/blob"
	"io/ioutil"
	_ "log"
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

func (t *TransformTool) Run() error {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")

	var region = flag.String("region", "full", "")
	var size = flag.String("size", "full", "")
	var rotation = flag.String("rotation", "0", "")
	var quality = flag.String("quality", "default", "")
	var format = flag.String("format", "jpg", "")

	var source_path = flag.String("source", "file:///", "...")
	var target_path = flag.String("target", "file:///", "...")

	var mode = flag.String("mode", "cli", "...")

	flag.Parse()

	// TO DO: validate args...

	// args := flag.Args()

	ctx := context.Background()

	source_bucket, err := bucket.OpenBucket(ctx, *source_path)

	if err != nil {
		return err
	}

	target_bucket, err := bucket.OpenBucket(ctx, *target_path)

	if err != nil {
		return err
	}

	if *cfg == "" {
		return errors.New("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	if err != nil {
		return err
	}

	driver, err := iiifdriver.NewDriverFromConfig(config)

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
