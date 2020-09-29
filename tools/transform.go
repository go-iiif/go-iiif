package tools

import (
	"context"
	"errors"
	"flag"
	aws_events "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiifdriver "github.com/go-iiif/go-iiif/v4/driver"
	iiifimage "github.com/go-iiif/go-iiif/v4/image"
	iiiflevel "github.com/go-iiif/go-iiif/v4/level"
	iiifsource "github.com/go-iiif/go-iiif/v4/source"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"gocloud.dev/blob"
	"io/ioutil"
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

	fh, err := opts.SourceBucket.NewReader(ctx, origin, nil)

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

	image, err := opts.Driver.NewImageFromConfigWithSource(opts.Config, source, origin)

	if err != nil {
		return err
	}

	err = image.Transform(opts.Transformation)

	if err != nil {
		return err
	}

	wr, err := opts.TargetBucket.NewWriter(ctx, target, nil)

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

func TransformToolFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flag.NewFlagSet("transform", flag.ExitOnError)

	err := AppendCommonTransformToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	err = AppendTransformToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func AppendCommonTransformToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	err := AppendCommonConfigFlags(ctx, fs)

	if err != nil {
		return err
	}

	err = AppendCommonToolModeFlags(ctx, fs)

	if err != nil {
		return err
	}

	return nil
}

func AppendTransformToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("region", "full", "A valid IIIF 2.0 region value.")
	fs.String("size", "full", "A valid IIIF 2.0 size value.")
	fs.String("rotation", "0", "A valid IIIF 2.0 rotation value.")
	fs.String("quality", "default", "A valid IIIF 2.0 quality value.")
	fs.String("format", "jpg", "A valid IIIF 2.0 format value.")

	fs.String("source", "file:///", "A valid Go Cloud bucket URI where the source file to transform is located.")
	fs.String("target", "file:///", "A valid Go Cloud bucket URI where the transformed file should be written.")

	return nil
}

func (t *TransformTool) Run(ctx context.Context) error {

	fs, err := TransformToolFlagSet(ctx)

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

func (t *TransformTool) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	paths := fs.Args()
	return t.RunWithFlagSetAndPaths(ctx, fs, paths...)
}

func (t *TransformTool) RunWithFlagSetAndPaths(ctx context.Context, fs *flag.FlagSet, paths ...string) error {

	config_name, err := lookup.StringVar(fs, "config-name")

	if err != nil {
		return err
	}

	config_source, err := lookup.StringVar(fs, "config-source")

	if err != nil {
		return err
	}

	region, err := lookup.StringVar(fs, "region")

	if err != nil {
		return err
	}

	size, err := lookup.StringVar(fs, "size")

	if err != nil {
		return err
	}

	rotation, err := lookup.StringVar(fs, "rotation")

	if err != nil {
		return err
	}

	quality, err := lookup.StringVar(fs, "quality")

	if err != nil {
		return err
	}

	format, err := lookup.StringVar(fs, "format")

	if err != nil {
		return err
	}

	source_path, err := lookup.StringVar(fs, "source")

	if err != nil {
		return err
	}

	target_path, err := lookup.StringVar(fs, "target")

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

	config_bucket, err := blob.OpenBucket(ctx, config_source)

	if err != nil {
		return err
	}

	config, err := iiifconfig.NewConfigFromBucket(ctx, config_bucket, config_name)

	if err != nil {
		return err
	}

	driver, err := iiifdriver.NewDriverFromConfig(config)

	if err != nil {
		return err
	}

	// TO DO DEFAULT TO source/target FROM config BUT CHECK FOR OVERRIDE IN *source/target_path ARGS

	source_bucket, err := blob.OpenBucket(ctx, source_path)

	if err != nil {
		return err
	}

	target_bucket, err := blob.OpenBucket(ctx, target_path)

	if err != nil {
		return err
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "http://127.0.0.1")

	if err != nil {
		return err
	}

	transformation, err := iiifimage.NewTransformation(level, region, size, rotation, quality, format)

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

	switch mode {

	case "cli":

		to_transform := make([]iiifuri.URI, 0)

		for _, str_uri := range paths {

			u, err := iiifuri.NewURI(str_uri)

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

		handler := func(ctx context.Context, ev aws_events.S3Event) error {

			to_transform := make([]iiifuri.URI, 0)

			for _, r := range ev.Records {

				s3_entity := r.S3
				s3_obj := s3_entity.Object
				s3_key := s3_obj.Key

				s3_fname := filepath.Base(s3_key)

				u, err := iiifuri.NewURI(s3_fname)

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
		return errors.New("Unsupported mode")
	}

	return nil
}
