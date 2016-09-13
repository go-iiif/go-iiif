package main

import (
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	iiifprofile "github.com/thisisaaronland/go-iiif/profile"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
	"github.com/whosonfirst/go-sanitize"
	"log"
	"net/http"
	"os"
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

	fs := http.FileServer(http.Dir(root))

	f := func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("GET", r.URL)
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

func ProfileHandlerFunc(config *iiifconfig.Config) (http.HandlerFunc, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		endpoint := EndpointFromRequest(r)
		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(level)

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

func InfoHandlerFunc(config *iiifconfig.Config) (http.HandlerFunc, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		opts := sanitize.DefaultOptions()

		vars := mux.Vars(r)

		id := vars["identifier"]
		id, _ = sanitize.SanitizeString(id, opts)

		id, err := iiifimage.ScrubIdentifier(id)

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
		profile, err := iiifprofile.NewProfile(endpoint, image)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

		rel_path := r.URL.Path

		/*

			for example:
			./bin/iiif-tile-seed -config config.json -scale-factors 1,2,4,8 -refresh 184512_5f7f47e5b3c66207_x.jpg

		*/

		body, err := derivatives_cache.Get(rel_path)

		if err == nil {

			cacheHit.Add(1)

			source, _ := iiifsource.NewMemorySource(body)
			image, _ := iiifimage.NewImageFromConfigWithSource(config, source, "cache")

			w.Header().Set("Content-Type", image.ContentType())
			w.Write(image.Body())
			return
		}

		opts := sanitize.DefaultOptions()
		vars := mux.Vars(r)

		id := vars["identifier"]
		id, _ = sanitize.SanitizeString(id, opts)

		region := vars["region"]
		region, _ = sanitize.SanitizeString(region, opts)

		size := vars["size"]
		size, _ = sanitize.SanitizeString(size, opts)

		rotation := vars["rotation"]
		rotation, _ = sanitize.SanitizeString(rotation, opts)

		quality := vars["quality"]
		quality, _ = sanitize.SanitizeString(quality, opts)

		format := vars["format"]
		format, _ = sanitize.SanitizeString(format, opts)

		endpoint := EndpointFromRequest(r)
		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transformation, err := iiifimage.NewTransformation(level, region, size, rotation, quality, format)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err = iiifimage.ScrubIdentifier(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		image, err := iiifimage.NewImageFromConfigWithCache(config, images_cache, id)

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

			}(rel_path, image)
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

	var host = flag.String("host", "localhost", "Define the hostname")
	var port = flag.Int("port", 8080, "Define which TCP port to use")
	var example = flag.Bool("example", false, "...")
	var root = flag.String("example-root", "example", "...")
	var cfg = flag.String("config", ".", "config")

	flag.Parse()

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

	ProfileHandler, err := ProfileHandlerFunc(config)

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

	router.HandleFunc("/level2.json", ProfileHandler)
	router.HandleFunc("/{identifier}/info.json", InfoHandler)
	router.HandleFunc("/{identifier}/{region}/{size}/{rotation}/{quality}.{format}", ImageHandler)

	expvarHandler, _ := ExpvarHandlerFunc(*host)
	router.HandleFunc("/debug/vars", expvarHandler)

	if *example {

		exampleHandler, err := ExampleHandler(*root)

		if err != nil {
			log.Fatal(err)
		}

		// router.HandleFunc("/example", exampleHandler)
		router.HandleFunc("/example/{ignore:.*}", exampleHandler)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	err = gracehttp.Serve(&http.Server{Addr: endpoint, Handler: router})

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
