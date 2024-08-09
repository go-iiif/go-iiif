package tools

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"github.com/aaronland/go-http-ping"
	"github.com/aaronland/go-http-server"
	iiifcache "github.com/go-iiif/go-iiif/v7/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
	iiifdriver "github.com/go-iiif/go-iiif/v7/driver"
	iiifhttp "github.com/go-iiif/go-iiif/v7/http"
	iiiflevel "github.com/go-iiif/go-iiif/v7/level"
	iiifsource "github.com/go-iiif/go-iiif/v7/source"
	"github.com/gorilla/mux"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
)

type IIIFServerTool struct {
	Tool
}

func NewIIIFServerTool() (Tool, error) {

	t := &IIIFServerTool{}
	return t, nil
}

func ServerToolFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flag.NewFlagSet("server", flag.ExitOnError)

	err := AppendCommonFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	err = AppendCommonServerToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	err = AppendServerToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func AppendCommonServerToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	err := AppendCommonConfigFlags(ctx, fs)

	if err != nil {
		return err
	}

	return nil
}

func AppendServerToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("protocol", "", "The protocol for iiif-server server to listen on. Valid protocols are: http, lambda. THIS FLAG IS DEPRECATED: Please use -server-uri instead.")
	fs.String("host", "", "Bind the server to this host. THIS FLAG IS DEPRECATED: Please use -server-uri instead.")
	fs.Int("port", 0, "Bind the server to this port. THIS FLAG IS DEPRECATED: Please use -server-uri instead.")

	fs.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI")

	fs.Bool("example", false, "Add an /example endpoint to the server for testing and demonstration purposes")
	fs.String("example-root", "example", "An explicit path to a folder containing example assets")

	return nil
}

func (t *IIIFServerTool) Run(ctx context.Context) error {

	fs, err := ServerToolFlagSet(ctx)

	if err != nil {
		return err
	}

	flagset.Parse(fs)

	err = flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return err
	}

	return t.RunWithFlagSet(ctx, fs)
}

func (t *IIIFServerTool) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	paths := fs.Args()
	return t.RunWithFlagSetAndPaths(ctx, fs, paths...)
}

func (t *IIIFServerTool) RunWithFlagSetAndPaths(ctx context.Context, fs *flag.FlagSet, paths ...string) error {

	proto, err := lookup.StringVar(fs, "protocol")

	if err != nil {
		return err
	}

	host, err := lookup.StringVar(fs, "host")

	if err != nil {
		return err
	}

	port, err := lookup.IntVar(fs, "port")

	if err != nil {
		return err
	}

	server_uri, err := lookup.StringVar(fs, "server-uri")

	if err != nil {
		return err
	}

	example, err := lookup.BoolVar(fs, "example")

	if err != nil {
		return err
	}

	example_root, err := lookup.StringVar(fs, "example-root")

	if err != nil {
		return err
	}

	config, err := iiifconfig.LoadConfigWithFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	driver, err := iiifdriver.NewDriverFromConfig(config)

	if err != nil {
		return err
	}

	/*
		See this - we're just going to make sure we have a valid source
		before we start serving images (20160901/thisisaaronland)
	*/

	_, err = iiifsource.NewSourceFromConfig(config)

	if err != nil {
		return err
	}

	_, err = iiiflevel.NewLevelFromConfig(config, host)

	if err != nil {
		return err
	}

	/*

		Okay now we're going to set up global cache thingies for source images
		and derivatives mostly to account for the fact that in-memory cache
		thingies need to be... well, global

	*/

	images_cache, err := iiifcache.NewImagesCacheFromConfig(config)

	if err != nil {
		return err
	}

	derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		return err
	}

	info_handler, err := iiifhttp.InfoHandler(config, driver)

	if err != nil {
		return err
	}

	image_handler, err := iiifhttp.ImageHandler(config, driver, images_cache, derivatives_cache)

	if err != nil {
		return err
	}

	ping_handler, err := ping.PingHandler()

	if err != nil {
		return err
	}

	expvar_handler, err := iiifhttp.ExpvarHandler(host)

	if err != nil {
		return err
	}

	router := mux.NewRouter()

	router.Handle("/ping", ping_handler)
	router.HandleFunc("/debug/vars", expvar_handler)

	// https://github.com/go-iiif/go-iiif/issues/4

	router.HandleFunc("/{identifier:.+}/info.json", info_handler)
	router.HandleFunc("/{identifier:.+}/{region}/{size}/{rotation}/{quality}.{format}", image_handler)

	if example {

		abs_path, err := filepath.Abs(example_root)

		if err != nil {
			return err
		}

		_, err = os.Stat(abs_path)

		if os.IsNotExist(err) {
			return err
		}

		example_handler, err := iiifhttp.ExampleHandler(abs_path)

		if err != nil {
			return err
		}

		router.HandleFunc("/example/{ignore:.*}", example_handler)
	}

	if proto != "" && host != "" && port != 0 {
		u := url.URL{}
		u.Scheme = proto
		u.Host = fmt.Sprintf("%s:%d", host, port)
		server_uri = u.String()
	}

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return err
	}

	slog.Info("Listening for requests", "address", s.Address())

	err = s.ListenAndServe(ctx, router)

	if err != nil {
		return err
	}

	return nil
}
