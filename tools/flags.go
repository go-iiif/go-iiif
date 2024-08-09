package tools

import (
	"context"
	"flag"
)

func AppendCommonFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.Bool("verbose", false, "Enabled verbose (debug) loggging.")
	return nil
}

func AppendCommonConfigFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("config-source", "", "A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used.")
	fs.String("config-name", "config.json", "The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'.")

	fs.String("config-images-source-uri", "", "If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.")
	fs.String("config-derivatives-cache-uri", "", "If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.")
	return nil
}

func AppendCommonInstructionsFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("instructions-source", "", "A valid Go Cloud bucket URI where your go-iiif \"instructions\" processing file is located. Optionally, if 'defaults://' is specified then the default instructions set bundled with this package will be used.")
	fs.String("instructions-name", "instructions.json", "The name of your go-iiif instructions file. This value will be ignored if -instructions-source is 'defaults://'.")

	return nil
}

func AppendCommonToolModeFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("mode", "cli", "Valid modes are: cli, csv, fsnotify, lambda.")
	return nil
}
