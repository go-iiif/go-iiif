package seed

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	"github.com/sfomuseum/go-flags/flagset"
)

type TileSeedOnCompleteFunc func(*iiifconfig.Config, string, string, int, error) error

type RunOptions struct {
	Config          *iiifconfig.Config        `json:"config"`
	Mode            string                    `json:"mode"`
	Endpoint        string                    `json:"endpoint"`
	Quality         string                    `json:"quality"`
	Format          string                    `json:"format"`
	Paths           []string                  `json:"paths"`
	ScaleFactors    []int                     `json:"scale_factors"`
	Refresh         bool                      `json:"refresh"`
	Workers         int                       `json:"workers"`
	NoExtension     bool                      `json:"no_extension"`
	GenerateHTML    bool                      `json:"generate_html"`
	OnCompleteFuncs []TileSeedOnCompleteFunc  `json:",omitempty"`
	URIFunc         iiifuri.URIInitializeFunc `json:",omitempty"`
	Verbose         bool                      `json:"verbose"`
}

func (o *RunOptions) AddOnCompleteFunc(fn TileSeedOnCompleteFunc) {
	o.OnCompleteFuncs = append(o.OnCompleteFuncs, fn)
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

	scales := make([]int, 0)

	for _, s := range strings.Split(scale_factors, ",") {

		s = strings.Trim(s, " ")
		scale, err := strconv.Atoi(s)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse scale factor, %w", err)
		}

		scales = append(scales, scale)
	}

	paths := fs.Args()

	opts := &RunOptions{
		Mode:         mode,
		Config:       cfg,
		Endpoint:     endpoint,
		Quality:      quality,
		Format:       format,
		ScaleFactors: scales,
		Paths:        paths,
		Workers:      processes,
		GenerateHTML: generate_html,
		Verbose:      verbose,
	}

	return opts, nil
}
