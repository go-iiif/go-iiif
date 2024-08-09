package server

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ConfigSource              string
	ConfigName                string
	ConfigImagesSourceURI     string
	ConfigDerivativesCacheURI string
	ServerURI                 string
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		ConfigSource:              config_source,
		ConfigName:                config_name,
		ConfigImagesSourceURI:     images_source_uri,
		ConfigDerivativesCacheURI: derivatives_cache_uri,
		ServerURI:                 server_uri,
	}

	return opts, nil
}
