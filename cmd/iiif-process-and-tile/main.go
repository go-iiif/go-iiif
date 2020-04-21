package main

import (
	_ "github.com/aaronland/go-cloud-s3blob"
	_ "github.com/go-iiif/go-iiif/v3/native"
	_ "gocloud.dev/blob/fileblob"
)

import (
	"context"
	"flag"
	"github.com/go-iiif/go-iiif/v3/tools"
	"github.com/sfomuseum/go-flags"
	"log"
)

func main() {

	ctx := context.Background()

	// set up flags for tools

	fs := flag.NewFlagSet("iiif-process-and-tile", flag.ExitOnError)

	err := tools.AppendCommonConfigFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append config flags, %v", err)
	}

	err = tools.AppendCommonInstructionsFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append config flags, %v", err)
	}

	err = tools.AppendCommonToolModeFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append tool flags, %v", err)
	}

	err = tools.AppendProcessToolFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append process tool flags, %v", err)
	}

	err = tools.AppendTileSeedToolFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append tileseed tool flags, %v", err)
	}

	// add custom flags

	fs.Bool("synchronous", false, "Run tools synchronously.")

	// parse flags

	flags.Parse(fs)

	err = flags.SetFlagsFromEnvVars(fs, "IIIF_PROCESS_AND_TILE")

	if err != nil {
		log.Fatalf("Failed to set flags from environment, %v", err)
	}

	// retrieve custom flags

	sync, err := flags.BoolVar(fs, "synchronous")

	if err != nil {
		log.Fatalf("Failed to parse -synchronous flag, %v", err)
	}

	// create tools

	pr_tool, err := tools.NewProcessTool()

	if err != nil {
		log.Fatalf("Failed to create new process tool, %v", err)
	}

	ts_tool, err := tools.NewTileSeedTool()

	if err != nil {
		log.Fatalf("Failed to create new process tool, %v", err)
	}

	// create tool runner

	var runner tools.Tool

	if sync {
		runner, err = tools.NewSynchronousToolRunner(pr_tool, ts_tool)
	} else {
		runner, err = tools.NewToolRunner(pr_tool, ts_tool)
	}

	if err != nil {
		log.Fatalf("Failed to create new tool runner, %v", err)
	}

	// run tools

	paths := fs.Args()
	err = runner.RunWithFlagSetAndPaths(ctx, fs, paths...)

	if err != nil {
		log.Fatalf("Failed to run process tool, %v", err)
	}

}
