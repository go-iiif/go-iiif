package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var config_source string
var config_name string
var images_source_uri string
var derivatives_cache_uri string

var server_uri string

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&config_source, "config-source", "defaults://", "...")
	fs.StringVar(&config_name, "config-name", "config.json", "...")
	fs.StringVar(&images_source_uri, "config-images-source-uri", "", "...")
	fs.StringVar(&derivatives_cache_uri, "config-derivatives-cache-uri", "mem://", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging")
	return fs
}
