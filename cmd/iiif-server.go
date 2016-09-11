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
	"log"
	"net/http"
	"os"
	"strings"
)

func ExpvarHandlerFunc(host string) (http.HandlerFunc, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		remote := strings.Split(r.RemoteAddr, ":")

		if remote[0] != host {

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

		level, err := iiiflevel.NewLevelFromConfig(config, r.Host)

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

		vars := mux.Vars(r)

		id, err := iiifimage.ScrubIdentifier(vars["identifier"])

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		image, err := iiifimage.NewImageFromConfig(config, id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		profile, err := iiifprofile.NewProfile(r.Host, image)

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

		body, err := derivatives_cache.Get(rel_path)

		if err == nil {

			// log.Println("CACHE HIT", rel_path)

			source, _ := iiifsource.NewMemorySource(body)
			image, _ := iiifimage.NewImageFromConfigWithSource(config, source, "cache")

			w.Header().Set("Content-Type", image.ContentType())
			w.Write(image.Body())
			return
		}

		vars := mux.Vars(r)
		id := vars["identifier"]

		region := vars["region"]
		size := vars["size"]
		rotation := vars["rotation"]
		quality := vars["quality"]
		format := vars["format"]

		level, err := iiiflevel.NewLevelFromConfig(config, r.Host)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transformation, err := iiifimage.NewTransformation(level, region, size, rotation, quality, format)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		image, err := iiifimage.NewImageFromConfigWithCache(config, images_cache, id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// something something something maybe sendfile something something
		// (20160901/thisisaaronland)
		//
		// log.Printf("%s %t\n", id, transformation.HasTransformation())

		if transformation.HasTransformation() {

			err = image.Transform(transformation)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			go func(k string, im iiifimage.Image) {

				derivatives_cache.Set(k, im.Body())

			}(rel_path, image)
		}

		w.Header().Set("Content-Type", image.ContentType())
		w.Write(image.Body())
		return
	}

	return http.HandlerFunc(f), nil
}

func main() {

	var host = flag.String("host", "localhost", "Define the hostname")
	var port = flag.Int("port", 8080, "Define which TCP port to use")
	var cfg = flag.String("config", ".", "config")

	flag.Parse()

	config, err := iiifconfig.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	// See this - we're just going to make sure we have a valid source
	// before we start serving images (20160901/thisisaaronland)

	_, err = iiifsource.NewSourceFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	_, err = iiiflevel.NewLevelFromConfig(config, *host)

	if err != nil {
		log.Fatal(err)
	}

	// Okay now we're going to set up global cache thingies for source images
	// and derivatives mostly to account for the fact that in-memory cache
	// thingies need to be... well, global

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

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	err = gracehttp.Serve(&http.Server{Addr: endpoint, Handler: router})

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}
