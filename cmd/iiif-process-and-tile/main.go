package main

import (
	_ "github.com/aaronland/go-cloud-s3blob"
	_ "github.com/go-iiif/go-iiif/native"
	_ "gocloud.dev/blob/fileblob"
)

import (
	"context"
	"github.com/go-iiif/go-iiif/tools"
	"github.com/sfomuseum/go-flags"
	"log"
	"flag"
)

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()
	
	fs := flag.NewFlagSet("combined", flag.ExitOnError)

	err := tools.AppendProcessToolFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append process tool flags, %v", err)
	}

	/*

> go run -mod vendor cmd/iiif-process-and-tile/main.go -h
combined flag redefined: config
panic: combined flag redefined: config

goroutine 1 [running]:
flag.(*FlagSet).Var(0xc000094360, 0x1d44ee0, 0xc0000a0490, 0x1a97955, 0x6, 0x1aca4a6, 0x5d)
	/usr/local/go/src/flag/flag.go:851 +0x4b8
flag.(*FlagSet).StringVar(...)
	/usr/local/go/src/flag/flag.go:754
flag.(*FlagSet).String(0xc000094360, 0x1a97955, 0x6, 0x0, 0x0, 0x1aca4a6, 0x5d, 0xc0000a0480)
	/usr/local/go/src/flag/flag.go:767 +0xa0
github.com/go-iiif/go-iiif/tools.AppendTileSeedToolFlags(0x1d4ffe0, 0xc0000ac380, 0xc000094360, 0x0, 0x0)
	/Users/asc/code/go-iiif/tools/tileseed.go:124 +0x61
main.main()
	/Users/asc/code/go-iiif/cmd/iiif-process-and-tile/main.go:32 +0x139
exit status 2

	*/

	/*
	err = tools.AppendTileSeedToolFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append tileseed tool flags, %v", err)
	}
	*/
	
	flags.Parse(fs)

	err = flags.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		log.Fatalf("Failed to set flags from environment, %v", err)
	}
	
	pr_tool, err := tools.NewProcessTool()

	if err != nil {
		log.Fatalf("Failed to create new process tool, %v", err)
	}

	err = pr_tool.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to run process tool, %v", err)
	}
}
