package tools

import (
	"flag"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiiflevel "github.com/go-iiif/go-iiif/level"
	_ "github.com/go-iiif/go-iiif/native"
	iiifsource "github.com/go-iiif/go-iiif/source"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type TransformTool struct {
	Tool
}

func NewTransformTool() (Tool, error) {

	t := &TransformTool{}
	return t, nil
}

func (t *TransformTool) Run() error {

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

	driver, err := iiifdriver.NewDriverFromConfig(config)

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

	image, err := driver.NewImageFromConfigWithSource(config, source, fname)

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

	return nil
}
