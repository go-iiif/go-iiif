package transform

/*

> go run -mod vendor cmd/iiif-transform/main.go -config-name config.json.example -config-source file:///usr/local/code/go-iiif/docs/ -region -1,-1,320,320 -quality sharpen -source file:///usr/local/code/go-iiif/example/images -target file:///usr/local/code/go-iiif/example/cache spanking-cat.jpg

*/

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"

	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifaws "github.com/go-iiif/go-iiif/v6/aws"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifimage "github.com/go-iiif/go-iiif/v6/image"
	iiiflevel "github.com/go-iiif/go-iiif/v6/level"
	iiifsource "github.com/go-iiif/go-iiif/v6/source"
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

	driver, err := iiifdriver.NewDriverFromConfig(opts.Config)

	if err != nil {
		return err
	}

	level, err := iiiflevel.NewLevelFromConfig(opts.Config, "http://127.0.0.1")

	if err != nil {
		return err
	}

	compliance := level.Compliance()

	transformation, err := iiifimage.NewTransformation(compliance, opts.Region, opts.Size, opts.Rotation, opts.Quality, opts.Format)

	if err != nil {
		return err
	}

	transform_opts := &TransformOptions{
		Config:         opts.Config,
		Driver:         driver,
		Transformation: transformation,
	}

	switch opts.Mode {

	case "cli":

		ctx := context.Background()
		to_transform := make([]iiifuri.URI, 0)

		for _, str_uri := range opts.Paths {

			u, err := iiifuri.NewURI(ctx, str_uri)

			if err != nil {
				return err
			}

			to_transform = append(to_transform, u)
		}

		err = TransformMany(ctx, transform_opts, to_transform...)

		if err != nil {
			return err
		}

	case "lambda":

		handler := func(ctx context.Context, ev iiifaws.Event) error {

			to_transform := make([]iiifuri.URI, 0)

			for _, r := range ev.Records {

				s3_entity := r.S3
				s3_obj := s3_entity.Object
				s3_key := s3_obj.Key

				s3_fname := filepath.Base(s3_key)

				u, err := iiifuri.NewURI(ctx, s3_fname)

				if err != nil {
					return err
				}

				to_transform = append(to_transform, u)
			}

			err := TransformMany(ctx, transform_opts, to_transform...)

			if err != nil {
				return err
			}

			return nil
		}

		aws_lambda.Start(handler)

	default:
		return fmt.Errorf("Unsupported mode")
	}

	return nil
}

type TransformOptions struct {
	Config         *iiifconfig.Config
	Driver         iiifdriver.Driver
	Transformation *iiifimage.Transformation
	SourceBucket   *blob.Bucket // DEPRECATED
	TargetBucket   *blob.Bucket // DEPRECATED
}

func TransformMany(ctx context.Context, opts *TransformOptions, uris ...iiifuri.URI) error {

	for _, uri := range uris {

		err := Transform(ctx, opts, uri)

		if err != nil {
			return err
		}
	}

	return nil
}

func Transform(ctx context.Context, opts *TransformOptions, uri iiifuri.URI) error {

	origin := uri.Origin()
	target, err := uri.Target(nil)

	if err != nil {
		return err
	}

	source_r, err := opts.SourceBucket.NewReader(ctx, origin, nil)

	if err != nil {
		return err
	}

	defer source_r.Close()

	if !opts.Transformation.HasTransformation() {
		return fmt.Errorf("No transformation")
	}

	// START OF FIX ME - load from config

	source_body, err := io.ReadAll(source_r)

	if err != nil {
		return err
	}

	source, err := iiifsource.NewMemorySource(source_body)

	if err != nil {
		return err
	}

	iiif_image, err := opts.Driver.NewImageFromConfigWithSource(opts.Config, source, origin)

	if err != nil {
		return err
	}

	// END OF FIX ME - load from config

	err = iiif_image.Transform(opts.Transformation)

	if err != nil {
		return err
	}

	// START OF FIX ME - load from config

	wr, err := opts.TargetBucket.NewWriter(ctx, target, nil)

	if err != nil {
		return err
	}

	_, err = wr.Write(iiif_image.Body())

	if err != nil {
		return err
	}

	err = wr.Close()

	if err != nil {
		return err
	}

	// END OF FIX ME - load from config

	return nil
}
