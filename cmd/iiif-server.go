package main

import (
	"encoding/json"
	"errors"
	"expvar"
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	iiifpalette "github.com/thisisaaronland/go-iiif/palette"
	iiifprofile "github.com/thisisaaronland/go-iiif/profile"
	iiifservice "github.com/thisisaaronland/go-iiif/service"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
	"github.com/whosonfirst/go-sanitize"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var cacheHit *expvar.Int
var cacheMiss *expvar.Int
var cacheSet *expvar.Int

var transformsCount *expvar.Int
var transformsAvgTime *expvar.Float

var transforms_counter int64
var transforms_timer int64

var timers_mu *sync.Mutex

type IIIFParameters struct {
	Identifier string
	Region     string
	Size       string
	Rotation   string
	Quality    string
	Format     string
}

type IIIFQueryParser struct {
	opts *sanitize.Options
	vars map[string]string
}

func NewIIIFQueryParser(r *http.Request) (*IIIFQueryParser, error) {

	opts := sanitize.DefaultOptions()
	vars := mux.Vars(r)

	p := IIIFQueryParser{
		opts: opts,
		vars: vars,
	}

	return &p, nil
}

func (p *IIIFQueryParser) GetIIIFParameter(key string) (string, error) {

	var err error

	value := p.vars[key]

	value, err = sanitize.SanitizeString(value, p.opts)

	if err != nil {
		return "", err
	}

	value, err = url.QueryUnescape(value)

	if err != nil {
		return "", err
	}

	// This should be already be stripped out by the time we get here but just
	// in case... (20160926/thisisaaronland)

	if strings.Contains(value, "../") {
		msg := fmt.Sprintf("Invalid key %s", key)
		err := errors.New(msg)
		return "", err
	}

	return value, nil
}

func (p *IIIFQueryParser) GetIIIFParameters() (*IIIFParameters, error) {

	id, err := p.GetIIIFParameter("identifier")

	if err != nil {
		return nil, err
	}

	region, err := p.GetIIIFParameter("region")

	if err != nil {
		return nil, err
	}

	size, err := p.GetIIIFParameter("size")

	if err != nil {
		return nil, err
	}

	rotation, err := p.GetIIIFParameter("rotation")

	if err != nil {
		return nil, err
	}

	quality, err := p.GetIIIFParameter("quality")

	if err != nil {
		return nil, err
	}

	format, err := p.GetIIIFParameter("format")

	if err != nil {
		return nil, err
	}

	params := IIIFParameters{
		Identifier: id,
		Region:     region,
		Size:       size,
		Rotation:   rotation,
		Quality:    quality,
		Format:     format,
	}

	return &params, nil
}

func init() {

	cacheHit = expvar.NewInt("CacheHit")
	cacheMiss = expvar.NewInt("CacheMiss")
	cacheSet = expvar.NewInt("CacheSet")

	transformsCount = expvar.NewInt("TransformsCount")
	transformsAvgTime = expvar.NewFloat("TransformsAvgTimeMS")

	transforms_counter = 0
	transforms_timer = 0

	timers_mu = new(sync.Mutex)
}

func ExampleHandler(root string) (http.HandlerFunc, error) {

	base := fmt.Sprintf("/%s", filepath.Base(root))

	fs := http.FileServer(http.Dir(root))

	f := func(w http.ResponseWriter, r *http.Request) {

		r.URL.Path = strings.Replace(r.URL.Path, base, "", 1)
		fs.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f), nil
}

func ExpvarHandlerFunc(host string) (http.HandlerFunc, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		remote := strings.Split(r.RemoteAddr, ":")

		if remote[0] != "127.0.0.1" && remote[0] != host {

			log.Printf("host '%s' remote '%s'\n", remote[0], host)
			http.Error(w, "No soup for you!", http.StatusForbidden)
			return
		}

		// This is copied wholesale from
		// https://golang.org/src/expvar/expvar.go

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "{\n")

		first := true

		expvar.Do(func(kv expvar.KeyValue) {
			if !first {
				fmt.Fprintf(w, ",\n")
			}

			first = false
			fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
		})

		fmt.Fprintf(w, "\n}\n")
	}

	return http.HandlerFunc(f), nil
}

func InfoHandlerFunc(config *iiifconfig.Config) (http.HandlerFunc, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		parser, err := NewIIIFQueryParser(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := parser.GetIIIFParameter("identifier")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		image, err := iiifimage.NewImageFromConfig(config, id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		endpoint := EndpointFromRequest(r)

		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		profile, err := iiifprofile.NewProfile(endpoint, image, level)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// start of colour palette stuff

		p, _ := iiifpalette.NewVibrantPalette()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		im, err := iiifimage.IIIFImageToGolangImage(image)

		c, err := p.Extract(im)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("COLOURS", c)

		s, _ := iiifservice.NewPaletteService(endpoint, c)

		profile.AddService(s)

		// end of colour palette stuff
		
		b, err := json.Marshal(profile)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(b)

	}

	return http.HandlerFunc(f), nil
}

func ImageHandlerFunc(config *iiifconfig.Config, images_cache iiifcache.Cache, derivatives_cache iiifcache.Cache) (http.HandlerFunc, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		/*

		   Okay, you see all of this? We're going to validate all the things including
		   a new transformation in order to be able to call ToURI() to account for the
		   fact that the "default" format is whatever the server wants it to be which
		   means we need to convert "default" into "color" or whatever in order for the
		   caching layer to work and find things generated by iiif-tile-seed. Good times...
		   (20160916/thisisaaronland)

		*/

		query, err := NewIIIFQueryParser(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		params, err := query.GetIIIFParameters()

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		endpoint := EndpointFromRequest(r)
		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transformation, err := iiifimage.NewTransformation(level, params.Region, params.Size, params.Rotation, params.Quality, params.Format)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uri, err := transformation.ToURI(params.Identifier)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body, err := derivatives_cache.Get(uri)

		if err == nil {

			cacheHit.Add(1)

			source, _ := iiifsource.NewMemorySource(body)
			image, _ := iiifimage.NewImageFromConfigWithSource(config, source, "cache")

			w.Header().Set("Content-Type", image.ContentType())
			w.Write(image.Body())
			return
		}

		image, err := iiifimage.NewImageFromConfigWithCache(config, images_cache, params.Identifier)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		/*
			something something something maybe sendfile something something
			(20160901/thisisaaronland)
		*/

		if transformation.HasTransformation() {

			cacheMiss.Add(1)

			t1 := time.Now()
			err = image.Transform(transformation)
			t2 := time.Since(t1)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			go func(t time.Duration) {

				ns := t.Nanoseconds()
				ms := ns / (int64(time.Millisecond) / int64(time.Nanosecond))

				timers_mu.Lock()

				counter := atomic.AddInt64(&transforms_counter, 1)
				timer := atomic.AddInt64(&transforms_timer, ms)

				avg := float64(timer) / float64(counter)

				transformsCount.Add(1)
				transformsAvgTime.Set(avg)

				timers_mu.Unlock()
			}(t2)

			go func(k string, im iiifimage.Image) {

				derivatives_cache.Set(k, im.Body())
				cacheSet.Add(1)

			}(uri, image)
		}

		w.Header().Set("Content-Type", image.ContentType())
		w.Write(image.Body())
		return
	}

	return http.HandlerFunc(f), nil
}

func EndpointFromRequest(r *http.Request) string {

	scheme := "http"

	if r.TLS != nil {
		scheme = "https"
	}

	endpoint := fmt.Sprintf("%s://%s", scheme, r.Host)
	return endpoint
}

func main() {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")
	var host = flag.String("host", "localhost", "Bind the server to this host")
	var port = flag.Int("port", 8080, "Bind the server to this port")
	var example = flag.Bool("example", false, "Add an /example endpoint to the server for testing and demonstration purposes")
	var root = flag.String("example-root", "example", "An explicit path to a folder containing example assets")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	/*

		See this - we're just going to make sure we have a valid source
		before we start serving images (20160901/thisisaaronland)

	*/

	_, err = iiifsource.NewSourceFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	_, err = iiiflevel.NewLevelFromConfig(config, *host)

	if err != nil {
		log.Fatal(err)
	}

	/*

		Okay now we're going to set up global cache thingies for source images
		and derivatives mostly to account for the fact that in-memory cache
		thingies need to be... well, global

	*/

	images_cache, err := iiifcache.NewImagesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	InfoHandler, err := InfoHandlerFunc(config)

	if err != nil {
		log.Fatal(err)
	}

	ImageHandler, err := ImageHandlerFunc(config, images_cache, derivatives_cache)

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// https://github.com/thisisaaronland/go-iiif/issues/4

	router.HandleFunc("/{identifier:.+}/info.json", InfoHandler)
	router.HandleFunc("/{identifier:.+}/{region}/{size}/{rotation}/{quality}.{format}", ImageHandler)

	expvarHandler, _ := ExpvarHandlerFunc(*host)
	router.HandleFunc("/debug/vars", expvarHandler)

	if *example {

		abs_path, err := filepath.Abs(*root)

		if err != nil {
			log.Fatal(err)
		}

		_, err = os.Stat(abs_path)

		if os.IsNotExist(err) {
			log.Fatal(err)
		}

		exampleHandler, err := ExampleHandler(abs_path)

		if err != nil {
			log.Fatal(err)
		}

		router.HandleFunc("/example/{ignore:.*}", exampleHandler)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	err = gracehttp.Serve(&http.Server{Addr: endpoint, Handler: router})

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
