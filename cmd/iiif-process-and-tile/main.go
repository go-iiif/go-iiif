package main

import (
	_ "github.com/aaronland/go-cloud-s3blob"
	_ "github.com/go-iiif/go-iiif/native"
	_ "gocloud.dev/blob/fileblob"
)

import (
	"context"
	"flag"
	"github.com/go-iiif/go-iiif/tools"
	"github.com/sfomuseum/go-flags"
	"log"
)

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	fs := flag.NewFlagSet("combined", flag.ExitOnError)

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

	flags.Parse(fs)

	err = flags.SetFlagsFromEnvVars(fs, "IIIF_PROCESS_AND_TILE")

	if err != nil {
		log.Fatalf("Failed to set flags from environment, %v", err)
	}

	pr_tool, err := tools.NewProcessTool()

	if err != nil {
		log.Fatalf("Failed to create new process tool, %v", err)
	}

	ts_tool, err := tools.NewTileSeedTool()

	if err != nil {
		log.Fatalf("Failed to create new process tool, %v", err)
	}

	runner, err := tools.NewToolRunner(pr_tool, ts_tool)

	if err != nil {
		log.Fatalf("Failed to create new combined tool, %v", err)
	}
	
	err = runner.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to run process tool, %v", err)
	}
	
}
