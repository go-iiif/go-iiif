package main

// ./bin/iiif-process -config config.json -instructions instructions.json -uri avocado.png
// {"avocado.png":{"b":"avocado.png/full/!2048,1536/0/color.jpg","d":"avocado.png/-1,-1,320,320/full/0/dither.jpg","o":"avocado.png/full/full/-1/color.jpg"}}

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-iiif/cache"	
	"github.com/thisisaaronland/go-iiif/config"
	"github.com/thisisaaronland/go-iiif/process"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
	"os"
)

func main() {

	var iiif_config = flag.String("config", "", "Path to a valid go-iiif config file.")
	var instructions = flag.String("instructions", "", "Path to a valid go-iiif processing instructions file.")

	var uris flags.MultiString
	flag.Var(&uris, "uri", "One or more valid IIIF URIs.")

	flag.Parse()

	instruction_set, err := process.ReadInstructions(*instructions)

	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewConfigFromFlag(*iiif_config)

	if err != nil {
		log.Fatal(err)
	}

	pr, err := process.NewParallelProcessor(cfg)

	if err != nil {
		log.Fatal(err)
	}

	dest_cache, err := cache.NewDerivativesCacheFromConfig(cfg)

	if err != nil {
		log.Fatal(err)
	}
	
	results := make(map[string]interface{})
	
	for _, uri := range uris {

		rsp, err := pr.ProcessURIWithInstructionSet(uri, instruction_set)

		if err != nil {
			log.Fatal(err)
		}

		enc_rsp, err := json.Marshal(rsp)

		if err != nil {
			log.Fatal(err)
		}

		key := fmt.Sprintf("%s.processed.json", uri)
		err = dest_cache.Set(key, enc_rsp)

		if err != nil {
			log.Fatal(err)
		}
		
		results[uri] = rsp
	}
		
	enc_results, err := json.Marshal(results)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(enc_results))
	os.Exit(0)
}
