package process

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifprocess "github.com/go-iiif/go-iiif/v6/process"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Config         *iiifconfig.Config             `json:"config"`
	Instructions   iiifprocess.IIIFInstructionSet `json:"instructions"`
	Mode           string                         `json:"mode"`
	Paths          []string                       `json:"paths"`
	Report         bool                           `json:"report"`
	ReportSource   string                         `json:"report_source"`
	ReportTemplate string                         `json:"report_template"`
	ReportHTML     bool                           `json:"report_html"`
	Verbose        bool                           `json:"verbose"`
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

	instructions, err := iiifprocess.LoadInstructions(ctx, instructions_source, instructions_name)

	if err != nil {
		return nil, err
	}

	// Instructions stuff here...

	paths := fs.Args()

	opts := &RunOptions{
		Mode:           mode,
		Config:         cfg,
		Instructions:   instructions,
		Report:         report,
		ReportSource:   report_source,
		ReportTemplate: report_template,
		Paths:          paths,
		Verbose:        verbose,
	}

	return opts, nil
}
