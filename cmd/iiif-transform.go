package main

// ./bin/iiif-transform -config config.json -quality dither -size ,300 /usr/local/cooperhewitt/iiif/images/184512_5f7f47e5b3c66207_x.jpg /vagrant/test2.jpg

/*

Important: This is still wet paint. It works so long as you use it in a very particular way.
Namely reading individual files from disk and writing them back to disk. It might stay that
way. It might grow the ability to load files from other sources. I'm not sure yet...
(20160927/thisisaaronland)

*/

import (
	"flag"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")

	var region = flag.String("region", "full", "")
	var size = flag.String("size", "full", "")
	var rotation = flag.String("rotation", "0", "")
	var quality = flag.String("quality", "default", "")
	var format = flag.String("format", "jpg", "")

	flag.Parse()

	// TO DO: validate args...

	args := flag.Args()
	infile := args[0]
	outfile := args[1]

	fname := filepath.Base(infile)

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "http://127.0.0.1")

	if err != nil {
		log.Fatal(err)
	}

	transformation, err := iiifimage.NewTransformation(level, *region, *size, *rotation, *quality, *format)

	if err != nil {
		log.Fatal(err)
	}

	// TO DO : compare extension of infile to 'format'

	if !transformation.HasTransformation() {
		log.Fatal("No transformation")
	}

	body, err := ioutil.ReadFile(infile)

	if err != nil {
		log.Fatal(err)
	}

	source, err := iiifsource.NewMemorySource(body)

	if err != nil {
		log.Fatal(err)
	}

	image, err := iiifimage.NewImageFromConfigWithSource(config, source, fname)

	if err != nil {
		log.Fatal(err)
	}

	err = image.Transform(transformation)

	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.Create(outfile)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()
	fh.Write(image.Body())
	fh.Sync()

	os.Exit(0)
}
