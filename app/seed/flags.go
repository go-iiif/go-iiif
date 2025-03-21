package seed

import (
	"flag"
	"runtime"

	"github.com/sfomuseum/go-flags/flagset"
)

var config_source string
var config_name string

var config_images_source_uri string
var config_derivatives_cache_uri string

var mode string
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

	fs.StringVar(&config_source, "config-source", "defaults://", "A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used.")

	fs.StringVar(&config_name, "config-name", "config.json", "The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'.")

	fs.StringVar(&config_images_source_uri, "config-images-source-uri", "", "If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.")

	fs.StringVar(&config_derivatives_cache_uri, "config-derivatives-cache-uri", "", "If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.")

	fs.StringVar(&mode, "mode", "cli", "Valid options are: cli, csv, fsnotify, lambda")

	fs.StringVar(&scale_factors, "scale-factors", "8,4,2,1", "A comma-separated list of scale factors to seed tiles with")
	fs.StringVar(&quality, "quality", "default", "A valid IIIF quality parameter - if \"default\" then the code will try to determine which format you've set as the default")
	fs.StringVar(&format, "format", "jpg", "A valid IIIF format parameter")

	fs.IntVar(&processes, "processes", runtime.NumCPU(), "The number of concurrent processes to use when tiling images")

	fs.BoolVar(&no_extension, "no-extension", false, "Remove any extension from destination folder name.")

	fs.StringVar(&endpoint, "endpoint", "http://localhost:8080", "The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image")

	fs.BoolVar(&generate_html, "generate-html", false, "If true then the tiles directory will be updated to include HTML/JavaScript/CSS assets to display tiles as a \"slippy\" map (using the leaflet-iiif.js library.")

	fs.BoolVar(&refresh, "refresh", false, "Refresh a tile even if already exists (default false)")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	return fs
}
