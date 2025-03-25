package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-server"
	iiifexample "github.com/go-iiif/go-iiif/v7/app/server/example"
	iiifcache "github.com/go-iiif/go-iiif/v7/cache"
	iiifdriver "github.com/go-iiif/go-iiif/v7/driver"
	iiifhttp "github.com/go-iiif/go-iiif/v7/http"
	iiiflevel "github.com/go-iiif/go-iiif/v7/level"
	iiifsource "github.com/go-iiif/go-iiif/v7/source"
)

func Run(ctx context.Context) error {

	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	driver, err := iiifdriver.NewDriver(ctx, opts.Config.Graphics.Driver)

	if err != nil {
		return fmt.Errorf("Failed to create new driver '%s', %w", opts.Config.Graphics.Driver, err)
	}

	_, err = iiifsource.NewSource(ctx, opts.Config.Images.Source.URI)

	if err != nil {
		return fmt.Errorf("Failed to create new source '%s', %w", opts.Config.Images.Source.URI, err)
	}

	host := "FIXME"

	_, err = iiiflevel.NewLevelFromConfig(opts.Config, opts.Host)

	if err != nil {
		return fmt.Errorf("Failed to create new level, %w", err)
	}

	/*

		Okay now we're going to set up global cache thingies for source images
		and derivatives mostly to account for the fact that in-memory cache
		thingies need to be... well, global

	*/

	images_cache, err := iiifcache.NewCache(ctx, opts.Config.Images.Cache.URI)

	if err != nil {
		return fmt.Errorf("Failed to create images cache, %w", err)
	}

	derivatives_cache, err := iiifcache.NewCache(ctx, opts.Config.Derivatives.Cache.URI)

	if err != nil {
		return fmt.Errorf("Failed to create derivatives cache, %w", err)
	}

	info_handler, err := iiifhttp.InfoHandler(opts.Config, driver, images_cache)

	if err != nil {
		return fmt.Errorf("Failed to create info handler, %w", err)
	}

	image_handler, err := iiifhttp.ImageHandler(opts.Config, driver, images_cache, derivatives_cache)

	if err != nil {
		return fmt.Errorf("Failed to create images handler, %w", err)
	}

	expvar_handler, err := iiifhttp.ExpvarHandler(host)

	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	if opts.Example {

		example_fs := http.FS(iiifexample.FS)
		example_handler := http.FileServer(example_fs)

		mux.Handle("/", example_handler)
	}

	mux.Handle("/debug/vars", expvar_handler)

	mux.Handle("/{identifier}/info.json", info_handler)
	mux.Handle("/{identifier}/{region}/{size}/{rotation}/{quality_dot_format}", image_handler)

	s, err := server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	slog.Info("Listening for requests", "address", s.Address())
	return s.ListenAndServe(ctx, mux)
}
