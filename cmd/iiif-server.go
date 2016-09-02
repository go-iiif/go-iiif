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

		l, err := iiiflevel.NewLevel2(r.Host)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(l)

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

func ImageHandlerFunc(config *iiifconfig.Config) (http.HandlerFunc, error) {

	cache, err := iiifcache.NewCacheFromConfig(config.Derivatives.Cache)

	if err != nil {
		return nil, err
	}

	f := func(w http.ResponseWriter, r *http.Request) {

		rel_path := r.URL.Path

		body, err := cache.Get(rel_path)

		if err == nil {

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

		transformation, err := iiifimage.NewTransformation(region, size, rotation, quality, format)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		image, err := iiifimage.NewImageFromConfig(config, id)

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

				cache.Set(k, im.Body())

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
		os.Exit(1)
	}

	// See this - we're just going to make sure we have a valid source
	// before we start serving images (20160901/thisisaaronland)

	_, err = iiifsource.NewSourceFromConfig(config.Images)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ProfileHandler, err := ProfileHandlerFunc(config)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	InfoHandler, err := InfoHandlerFunc(config)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ImageHandler, err := ImageHandlerFunc(config)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
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
