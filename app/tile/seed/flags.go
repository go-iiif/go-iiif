package seed

import (
	"flag"
	"runtime"

	"github.com/sfomuseum/go-flags/flagset"
)

var mode string
var csv_source string
var scale_factors string
var quality string
var format string
var processes int
var no_extension bool
var refresh bool
var endpoint string
var generate_html bool

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("seed")

	fs.StringVar(&mode, "mode", "cli", "Valid options are: cli, csv, lambda")
	fs.StringVar(&csv_source, "csv-source", "A valid Go Cloud bucket URI where your CSV tileseed files are located.", "")

	fs.StringVar(&scale_factors, "scale-factors", "4", "A comma-separated list of scale factors to seed tiles with")
	fs.StringVar(&quality, "quality", "default", "A valid IIIF quality parameter - if \"default\" then the code will try to determine which format you've set as the default")
	fs.StringVar(&format, "format", "jpg", "A valid IIIF format parameter")

	fs.IntVar(&processes, "processes", runtime.NumCPU(), "The number of concurrent processes to use when tiling images")

	fs.BoolVar(&no_extension, "noextension", false, "Remove any extension from destination folder name.")

	fs.BoolVar(&refresh, "refresh", false, "Refresh a tile even if already exists (default false)")
	fs.StringVar(&endpoint, "endpoint", "http://localhost:8080", "The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image")

	fs.BoolVar(&generate_html, "generate-tiles-html", false, "If true then the tiles directory will be updated to include HTML/JavaScript/CSS assets to display tiles as a \"slippy\" map (using the leaflet-iiif.js library.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return nil
}

/*

func TileSeedToolFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flag.NewFlagSet("tileseed", flag.ExitOnError)

	err := AppendCommonTileSeedToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	err = AppendTileSeedToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func AppendCommonTileSeedToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	err := AppendCommonFlags(ctx, fs)

	if err != nil {
		return err
	}

	err = AppendCommonConfigFlags(ctx, fs)

	if err != nil {
		return err
	}

	err = AppendCommonToolModeFlags(ctx, fs)

	if err != nil {
		return err
	}

	return nil
}

*/
