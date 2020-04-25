package tools

import (
	"context"
	"errors"
	"flag"
	"fmt"
	iiifcache "github.com/go-iiif/go-iiif/v4/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiifdriver "github.com/go-iiif/go-iiif/v4/driver"
	iiifhttp "github.com/go-iiif/go-iiif/v4/http"
	iiiflevel "github.com/go-iiif/go-iiif/v4/level"
	iiifserver "github.com/go-iiif/go-iiif/v4/server"
	iiifsource "github.com/go-iiif/go-iiif/v4/source"
	"github.com/gorilla/mux"
	"github.com/sfomuseum/go-flags"
	"gocloud.dev/blob"
	"log"
	"net/url"
	"os"
	"path/filepath"
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

	err := AppendCommonServerToolFlags(ctx, fs)

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

	fs.String("protocol", "http", "The protocol for wof-staticd server to listen on. Valid protocols are: http, lambda.")
	fs.String("host", "localhost", "Bind the server to this host")
	fs.Int("port", 8080, "Bind the server to this port")

	fs.Bool("example", false, "Add an /example endpoint to the server for testing and demonstration purposes")
	fs.String("example-root", "example", "An explicit path to a folder containing example assets")

	return nil
}

func (t *IIIFServerTool) Run(ctx context.Context) error {

	fs, err := ServerToolFlagSet(ctx)

	if err != nil {
		return err
	}

	flags.Parse(fs)

	err = flags.SetFlagsFromEnvVars(fs, "IIIF_SERVER")

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

	cfg, err := flags.StringVar(fs, "config")

	if err != nil {
		return err
	}

	config_source, err := flags.StringVar(fs, "config-source")

	if err != nil {
		return err
	}

	config_name, err := flags.StringVar(fs, "config-name")

	if err != nil {
		return err
	}

	proto, err := flags.StringVar(fs, "protocol")

	if err != nil {
		return err
	}

	host, err := flags.StringVar(fs, "host")

	if err != nil {
		return err
	}

	port, err := flags.IntVar(fs, "port")

	if err != nil {
		return err
	}

	example, err := flags.BoolVar(fs, "example")

	if err != nil {
		return err
	}

	example_root, err := flags.StringVar(fs, "example-root")

	if err != nil {
		return err
	}

	if cfg != "" {

		log.Println("-config flag is deprecated. Please use -config-source and -config-name (setting them now).")

		abs_config, err := filepath.Abs(cfg)

		if err != nil {
			return err
		}

		config_name = filepath.Base(abs_config)
		config_source = fmt.Sprintf("file://%s", filepath.Dir(abs_config))
	}

	if config_source == "" {
		return errors.New("Required -config-source flag is empty.")
	}

	config_bucket, err := blob.OpenBucket(ctx, config_source)

	if err != nil {
		return err
	}

	config, err := iiifconfig.NewConfigFromBucket(ctx, config_bucket, config_name)

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

	ping_handler, err := iiifhttp.PingHandler()

	if err != nil {
		return err
	}

	expvar_handler, err := iiifhttp.ExpvarHandler(host)

	if err != nil {
		return err
	}

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping_handler)
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

	address := fmt.Sprintf("http://%s:%d", host, port)

	u, err := url.Parse(address)

	if err != nil {
		return err
	}

	s, err := iiifserver.NewServer(proto, u)

	if err != nil {
		return err
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(router)

	if err != nil {
		return err
	}

	return nil
}
