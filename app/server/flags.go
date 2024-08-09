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

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	return fs
}
