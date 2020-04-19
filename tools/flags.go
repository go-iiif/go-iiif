package tools

import (
	"context"
	"flag"
)

func AppendCommonConfigFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("config", "", "Path to a valid go-iiif config file. DEPRECATED - please use -config_source and -config name.")

	fs.String("config-source", "", "A valid Go Cloud bucket URI where your go-iiif config file is located.")
	fs.String("config-name", "config.json", "The name of your go-iiif config file.")

	return nil
}

func AppendCommonInstructionsFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("instructions", "", "Path to a valid go-iiif processing instructions file. DEPRECATED - please use -instructions-source and -instructions-name.")
	fs.String("instructions-source", "", "A valid Go Cloud bucket URI where your go-iiif instructions file is located.")
	fs.String("instructions-name", "instructions.json", "The name of your go-iiif instructions file.")

	return nil
}

func AppendCommonToolModeFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("mode", "cli", "Valid modes are: cli, csv, fsnotify, lambda.")
	return nil
}
