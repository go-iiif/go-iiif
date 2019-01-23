package main

/*

./bin/iiif-process -config ./config.json -instructions file://instructions.json example.jpg
2019/01/22 16:13:07 d example.jpg/-1,-1,320,320/full/0/dither.jpg
2019/01/22 16:13:07 b example.jpg/full/!2048,1536/0/color.jpg
2019/01/22 16:13:07 o example.jpg/full/full/-1/color.jpg

*/

import (
	"flag"
	// "github.com/thisisaaronland/go-iiif/cache"
	"github.com/thisisaaronland/go-iiif/config"
	// "github.com/thisisaaronland/go-iiif/image"
	"github.com/thisisaaronland/go-iiif/process"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
)

func main() {

	var iiif_config = flag.String("config", "", "Path to a valid go-iiif config file")
	var instructions = flag.String("instructions", "", "...")

	var uris flags.MultiString
	flag.Var(&uris, "uri", "...")

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

	/*
		dest_cache, err := cache.NewDerivativesCacheFromConfig(cfg)

		if err != nil {
			log.Fatal(err)
		}
	*/

	for _, uri := range uris {

		r, err := pr.ProcessURIWithInstructionSet(uri, instruction_set)

		if err != nil {
			log.Fatal(err)
		}

		for label, new_uri := range r {

			log.Println(label, new_uri)

			/*
				im, err := image.NewImageFromConfigWithCache(cfg, dest_cache, new_uri)

				if err != nil {
					log.Fatal(err)
				}
			*/
		}

	}

}
