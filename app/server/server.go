package server

import (
	"context"
	"flag"
	_ "fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-server"
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifhttp "github.com/go-iiif/go-iiif/v6/http"
	iiiflevel "github.com/go-iiif/go-iiif/v6/level"
	iiifsource "github.com/go-iiif/go-iiif/v6/source"
	iiifexample "github.com/go-iiif/go-iiif/v6/static/example"
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

	driver, err := iiifdriver.NewDriverFromConfig(opts.Config)

	if err != nil {
		return err
	}

	/*
		See this - we're just going to make sure we have a valid source
		before we start serving images (20160901/thisisaaronland)
	*/

	_, err = iiifsource.NewSourceFromConfig(opts.Config)

	if err != nil {
		return err
	}

	host := "FIXME"

	_, err = iiiflevel.NewLevelFromConfig(opts.Config, host)

	if err != nil {
		return err
	}

	/*

		Okay now we're going to set up global cache thingies for source images
		and derivatives mostly to account for the fact that in-memory cache
		thingies need to be... well, global

	*/

	images_cache, err := iiifcache.NewImagesCacheFromConfig(opts.Config)

	if err != nil {
		return err
	}

	derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(opts.Config)

	if err != nil {
		return err
	}

	info_handler, err := iiifhttp.InfoHandler(opts.Config, driver)

	if err != nil {
		return err
	}

	image_handler, err := iiifhttp.ImageHandler(opts.Config, driver, images_cache, derivatives_cache)

	if err != nil {
		return err
	}

	expvar_handler, err := iiifhttp.ExpvarHandler(host)

	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/debug/vars", expvar_handler)

	// https://github.com/go-iiif/go-iiif/issues/4

	mux.Handle("/{identifier:.+}/info.json", info_handler)
	mux.Handle("/{identifier:.+}/{region}/{size}/{rotation}/{quality}.{format}", image_handler)

	if opts.Example {

		http_fs := http.FS(iiifexample.FS)
		example_handler := http.FileServer(http_fs)

		mux.Handle("/example/", example_handler)
	}

	s, err := server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return err
	}

	slog.Info("Listening for requests", "address", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return err
	}

	return nil
}
