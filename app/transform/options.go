package transform

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Config   *iiifconfig.Config `json:"config"`
	Mode     string             `json:"mode"`
	Region   string             `json:"region"`
	Size     string             `json:"size"`
	Rotation string             `json:"rotation"`
	Quality  string             `json:"quality"`
	Format   string             `json:"format"`
	Paths    []string           `json:"paths"`
	Verbose  bool               `json:"verbose"`
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign tileseed tool flags from environment variables, %w", err)
	}

	cfg, err := iiifconfig.LoadConfig(ctx, config_source, config_name)

	if err != nil {
		return nil, err
	}

	if config_images_source_uri != "" {
		slog.Debug("Reassign images source", "uri", config_images_source_uri)
		cfg.Images.Source.URI = config_images_source_uri
	}

	if config_derivatives_cache_uri != "" {
		slog.Debug("Reassign derivatives cache", "uri", config_derivatives_cache_uri)
		cfg.Derivatives.Cache.URI = config_derivatives_cache_uri
	}

	paths := fs.Args()

	opts := &RunOptions{
		Mode:     mode,
		Config:   cfg,
		Region:   region,
		Size:     size,
		Rotation: rotation,
		Quality:  quality,
		Format:   format,
		Paths:    paths,
		Verbose:  verbose,
	}

	return opts, nil
}
