package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-server"
	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
	iiifdriver "github.com/go-iiif/go-iiif/v7/driver"
	iiifhttp "github.com/go-iiif/go-iiif/v7/http"
	"github.com/rs/cors"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive run options from flag set, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	// Set up config and driver
	
	cfg, err := iiifconfig.LoadConfig(ctx, opts.ConfigSource, opts.ConfigName)

	if err != nil {
		return fmt.Errorf("Failed to load config, %w", err)
	}

	if opts.ConfigImagesSourceURI != "" {
		slog.Debug("Reassign images source", "uri", opts.ConfigImagesSourceURI)
		cfg.Images.Source.URI = opts.ConfigImagesSourceURI
	}

	if opts.ConfigDerivativesCacheURI != "" {
		slog.Debug("Reassign derivatives cache", "uri", opts.ConfigDerivativesCacheURI)
		cfg.Derivatives.Cache.URI = opts.ConfigDerivativesCacheURI
	}

	driver, err := iiifdriver.NewDriverFromConfig(cfg)

	if err != nil {
		return fmt.Errorf("Failed to load driver, %w", err)
	}

	// Set up HTTP mux
	
	mux := http.NewServeMux()

	// Info handler

	info_handler, err := iiifhttp.InfoHandler(cfg, driver)

	if err != nil {
		return fmt.Errorf("Failed to create info handler, %w", err)
	}

	// This makes the compiler sad:
	// cors_wrapper := cors.Default()
	// info_handler = cors_wrapper.Handler(info_handler)

	// This does not because... computers:
	mux.Handle("/{identifier}/info.json", cors.Default().Handler(info_handler))

	// Image handler
	
	image_handler, err := iiifhttp.ImageHandler(cfg, driver)

	if err != nil {
		return fmt.Errorf("Failed to create image handler, %w", err)
	}

	mux.Handle("/{identifier}/{region}/{size}/{rotation}/{quality_dot_format}", image_handler)

	// Expvar handler
	
	expvar_handler, err := iiifhttp.ExpvarHandler("127.0.0.1")

	if err != nil {
		return fmt.Errorf("Failed to create expvar handler, %w", err)
	}

	mux.Handle("/debug/expvar", expvar_handler)

	// Start the server
	
	s, err := server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	slog.Info("Listening for requests", "address", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}
