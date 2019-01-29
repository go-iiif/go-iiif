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
	"path/filepath"
	"sync"
)

func main() {

	var iiif_config = flag.String("config", "", "Path to a valid go-iiif config file.")
	var instructions = flag.String("instructions", "", "Path to a valid go-iiif processing instructions file.")

	var report = flag.Bool("report", false, "Store a process report (JSON) for each URI in the cache tree.")
	var report_name = flag.String("report-name", "process.json", "The filename for process reports. Default is 'process.json' as in '${URI}/process.json'.")

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

	pr, err := process.NewIIIFProcessor(cfg)

	if err != nil {
		log.Fatal(err)
	}

	results := make(map[string]interface{})
	wg := new(sync.WaitGroup)

	for _, uri := range uris {

		u, err := process.NewIIIFURI(uri)

		if err != nil {
			log.Fatal(err)
		}

		rsp, err := process.ParallelProcessURIWithInstructionSet(cfg, pr, instruction_set, u)

		if err != nil {
			log.Fatal(err)
		}

		if *report {

			key := filepath.Join(uri, *report_name)
			wg.Add(1)

			go func() {

				defer wg.Done()
				err := report_processing(cfg, key, rsp)

				if err != nil {
					log.Printf("Unable to write process report %s, %s", key, err)
				}
			}()
		}

		results[uri] = rsp
	}

	wg.Wait()

	enc_results, err := json.Marshal(results)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(enc_results))
	os.Exit(0)
}

func report_processing(cfg *config.Config, key string, rsp map[string]interface{}) error {

	dest_cache, err := cache.NewDerivativesCacheFromConfig(cfg)

	if err != nil {
		return err

	}

	enc_rsp, err := json.Marshal(rsp)

	if err != nil {
		return err
	}

	return dest_cache.Set(key, enc_rsp)
}
