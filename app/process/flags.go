package process

import (
	"flag"

	iiifprocess "github.com/go-iiif/go-iiif/v6/process"
	"github.com/sfomuseum/go-flags/flagset"
)

var config_source string
var config_name string

var config_images_source_uri string
var config_derivatives_cache_uri string

var instructions_source string
var instructions_name string

var mode string
var verbose bool

var report bool
var report_source string
var report_template string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&mode, "mode", "cli", "Valid options are: cli, lambda")
	fs.BoolVar(&verbose, "verbose", false, "Enabled verbose (debug) loggging.")

	fs.StringVar(&config_source, "config-source", "defaults://", "A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used.")

	fs.StringVar(&config_name, "config-name", "config.json", "The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'.")

	fs.StringVar(&config_images_source_uri, "config-images-source-uri", "", "If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.")
	fs.StringVar(&config_derivatives_cache_uri, "config-derivatives-cache-uri", "", "If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.")

	fs.StringVar(&instructions_source, "instructions-source", "", "A valid Go Cloud bucket URI where your go-iiif \"instructions\" processing file is located. Optionally, if 'defaults://' is specified then the default instructions set bundled with this package will be used.")
	fs.StringVar(&instructions_name, "instructions-name", "instructions.json", "The name of your go-iiif instructions file. This value will be ignored if -instructions-source is 'defaults://'.")

	fs.BoolVar(&report, "report", false, "Store a process report (JSON) for each URI in the cache tree.")
	fs.StringVar(&report_template, "report-template", iiifprocess.REPORTNAME_TEMPLATE, "A valid URI template for generating process report filenames.")
	fs.StringVar(&report_source, "report-source", "", "A valid Go Cloud bucket URI where your report file will be saved. If empty reports will be stored alongside derivative (or cached) images.")

	// fs.Bool("generate-report-html", false, "Generate an HTML page showing all the images listed in a process report.")

	return fs
}
