package server

import (
	"flag"
	"fmt"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Config    *iiifconfig.Config `json:"config"`
	ServerURI string             `json:"server_uri"`
	Example   bool               `json:"example"`
	Verbose   bool               `json:"verbose"`
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign tileseed tool flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		ServerURI: server_uri,
		Verbose:   verbose,
	}

	return opts, nil
}
