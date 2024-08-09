package main

/*

	$> ./bin/iiif-process-and-tile \
		-config-source defaults:// \
		-instructions-source defaults:// \
		-verbose \
		-scale-factors '1,2,4,8' \
		-noextension \
		-report \
		-refresh \
		-config-images-source-uri file:///usr/local \
		-config-derivatives-cache-uri file:///usr/local/test \
		-tiles-prefix tiles \
		-generate-tiles-html \
		'idsecret:///IMG_9998.jpg?id=9998&secret=abc&secret_o=def&format=jpg&label=x'

*/

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"

	_ "github.com/aaronland/gocloud-blob/s3"
	"github.com/go-iiif/go-iiif-uri"
	_ "github.com/go-iiif/go-iiif/v7/native"
	"github.com/go-iiif/go-iiif/v7/tools"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
)

func main() {

	ctx := context.Background()

	// set up flags for tools

	fs := flag.NewFlagSet("iiif-process-and-tile", flag.ExitOnError)

	err := tools.AppendCommonFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append config flags, %v", err)
	}

	err = tools.AppendCommonConfigFlags(ctx, fs)

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

	fs.String("tiles-prefix", "", "A relative URL to use a prefix when storing tiles.")

	// parse flags

	flagset.Parse(fs)

	err = flagset.SetFlagsFromEnvVars(fs, "IIIF_PROCESS_AND_TILE")

	if err != nil {
		log.Fatalf("Failed to set flags from environment, %v", err)
	}

	// retrieve custom lookup

	sync, err := lookup.BoolVar(fs, "synchronous")

	if err != nil {
		log.Fatalf("Failed to parse -synchronous flag, %v", err)
	}

	tiles_prefix, err := lookup.StringVar(fs, "tiles-prefix")

	if err != nil {
		log.Fatalf("Failed to parse -tiles-prefix flag, %v", err)
	}

	// create tools

	pr_tool, err := tools.NewProcessTool()

	if err != nil {
		log.Fatalf("Failed to create new process tool, %v", err)
	}

	var ts_tool tools.Tool

	if tiles_prefix != "" {

		ts_uri_func := func(raw_uri string) (uri.URI, error) {

			slog.Debug("Derive tileseed URI func for prefix", "uri", raw_uri, "prefix", tiles_prefix)

			/*

				what the following code suggests is that the go-iiif-uri.URI
				interface needs to be updated to return query parameters
				(20200929/thisisaaronland)

			*/

			u, err := uri.NewURI(ctx, raw_uri)

			if err != nil {
				return nil, fmt.Errorf("Failed to derive new URI from '%s', %w", raw_uri, err)
			}

			u2, err := url.Parse(raw_uri)

			if err != nil {
				return nil, fmt.Errorf("Failed parse URI '%s', %w", raw_uri, err)
			}

			q := u2.Query()

			target, err := u.Target(&q)

			if err != nil {
				return nil, fmt.Errorf("Failed to derive target from URI %s, %w", u, err)
			}

			origin := u.Origin()

			if !strings.HasPrefix(origin, "/") {
				origin = fmt.Sprintf("/%s", origin)
			}

			root := filepath.Dir(target)
			path := filepath.Join(root, "tiles")

			file_uri := fmt.Sprintf("%s://%s?target=%s", uri.FILE_SCHEME, origin, path)

			slog.Debug("Tile seed URI", "uri", file_uri)
			return uri.NewFileURI(ctx, file_uri)
		}

		t, err := tools.NewTileSeedToolWithURIFunc(ts_uri_func)

		if err != nil {
			log.Fatalf("Failed to create new process tool, %v", err)
		}

		ts_tool = t

	} else {

		t, err := tools.NewTileSeedTool()

		if err != nil {
			log.Fatalf("Failed to create new process tool, %v", err)
		}

		ts_tool = t
	}

	// create tool runner

	opts := &tools.ToolRunnerOptions{
		Tools: []tools.Tool{
			pr_tool,
			ts_tool,
		},
	}

	if sync {

		throttle, err := tools.NewToolRunnerThrottle()

		if err != nil {
			log.Fatalf("Failed to create new tool runner throttle, %v", err)
		}

		opts.Throttle = throttle
	}

	opts.OnCompleteFunc = func(ctx context.Context, path string) error {
		log.Printf("Finished processing %s\n", path)
		return nil
	}

	runner, err := tools.NewToolRunnerWithOptions(opts)

	if err != nil {
		log.Fatalf("Failed to create new tool runner, %v", err)
	}

	paths := fs.Args()
	err = runner.RunWithFlagSetAndPaths(ctx, fs, paths...)

	if err != nil {
		log.Fatalf("Failed to run process tool, %v", err)
	}

}
