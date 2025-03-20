package process

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var config_source string
var config_name string

var config_images_source_uri string
var config_derivatives_cache_uri string

var instructions_source string
var instructions_name string

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.BoolVar(&verbose, "verbose", false, "Enabled verbose (debug) loggging.")

	fs.StringVar(&config_source, "config-source", "defaults://", "A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used.")

	fs.StringVar(&config_name, "config-name", "config.json", "The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'.")

	fs.StringVar(&config_images_source_uri, "config-images-source-uri", "", "If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.")
	fs.StringVar(&config_derivatives_cache_uri, "config-derivatives-cache-uri", "", "If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.")

	fs.StringVar(&instructions_source, "instructions-source", "", "A valid Go Cloud bucket URI where your go-iiif \"instructions\" processing file is located. Optionally, if 'defaults://' is specified then the default instructions set bundled with this package will be used.")
	fs.StringVar(&instructions_name, "instructions-name", "instructions.json", "The name of your go-iiif instructions file. This value will be ignored if -instructions-source is 'defaults://'.")

	return fs
}
