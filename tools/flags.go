package tools

import (
	"context"
	"flag"
)

func AppendCommonConfigFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("config-source", "", "A valid Go Cloud bucket URI where your go-iiif config file is located.")
	fs.String("config-name", "config.json", "The name of your go-iiif config file.")

	fs.String("config-images-source-uri", "", "...")
	fs.String("config-derivatives-cache-uri", "", "...")	
	return nil
}

func AppendCommonInstructionsFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("instructions-source", "", "A valid Go Cloud bucket URI where your go-iiif instructions file is located.")
	fs.String("instructions-name", "instructions.json", "The name of your go-iiif instructions file.")

	return nil
}

func AppendCommonToolModeFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("mode", "cli", "Valid modes are: cli, csv, fsnotify, lambda.")
	return nil
}
