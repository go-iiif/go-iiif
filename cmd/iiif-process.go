package main

import (
	"encoding/json"
	"flag"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"	
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type IIIFInstructionSet []IIIFInstructions

type IIIFInstructions struct {
	Identifier string
	Region     string
	Size       string
	Rotation   string
	Quality    string
	Format     string
}

func main() {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")

	var instructions = flag.String("instructions-set", "", "...")
	
	flag.Parse()
	
	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "http://localhost")

	if err != nil {
		log.Fatal(err)
	}
	
	images_cache, err := iiifcache.NewImagesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	var raw_instructions []byte

	if strings.HasPrefix(*instructions, "file://"){

		fh, err := os.Open(*instructions)

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			log.Fatal(err)
		}

		raw_instructions = body
	} else {

		raw_instructions = []byte(*instructions)
	}
	
	var instruction_set IIIFInstructionSet

	err = json.Unmarshal(raw_instructions, &instruction_set)

	if err != nil {
		log.Fatal(err)
	}
	
	for _, uri := range flag.Args() {

		err := ProcessURI(uri, instruction_set, config, level, images_cache, derivatives_cache)

		if err != nil {
			log.Println(uri, err)
		}
	}

}

func ProcessURI(uri string, instruction_set IIIFInstructionSet, config *iiifconfig.Config, level iiiflevel.Level, images_cache iiifcache.Cache, derivatives_cache iiifcache.Cache) error {

	done_ch := make(chan bool)
	err_ch := make(chan error)
		
	remaining := len(instruction_set)

	for _, i := range instruction_set {

		go func(uri string, i IIIFInstructions, config *iiifconfig.Config, level iiiflevel.Level, images_cache iiifcache.Cache, derivatives_cache iiifcache.Cache){

			defer func(){
				done_ch <- true
			}()

			new_uri, err := ProcessURIWithInstructions(uri, i, config, level, images_cache, derivatives_cache)
				
			if err != nil {
				err_ch <- err
				return
			}

			log.Println(uri, new_uri)
			
		}(uri, i, config, level, images_cache, derivatives_cache)
	}

	for {

		select {
		case <- done_ch:
			remaining -= 1
		case e := <- err_ch:
			log.Println(e)
		default:
			//
		}
	}

	return nil
}

func ProcessURIWithInstructions(uri string, i IIIFInstructions, config *iiifconfig.Config, level iiiflevel.Level, images_cache iiifcache.Cache, derivatives_cache iiifcache.Cache) (string, error) {

	image, err := iiifimage.NewImageFromConfigWithCache(config, images_cache, uri)

	if err != nil {
		return "", err
	}

	transformation, err := iiifimage.NewTransformation(level, i.Region, i.Size, i.Rotation, i.Quality, i.Format)

	if err != nil {
		return "", err
	}

	new_uri, err := transformation.ToURI(uri)

	if err != nil {
		return "", err
	}
	
	err = image.Transform(transformation)

	if err != nil {
		return "", err
	}
	
	derivatives_cache.Set(new_uri, image.Body())

	return new_uri, nil
}
	
