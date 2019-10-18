package tools

import (
	"errors"
	"flag"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiiflevel "github.com/go-iiif/go-iiif/level"
	iiifsource "github.com/go-iiif/go-iiif/source"
	"io/ioutil"
	_ "log"
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
		return errors.New("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	if err != nil {
		return err
	}

	driver, err := iiifdriver.NewDriverFromConfig(config)

	if err != nil {
		return err
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "http://127.0.0.1")

	if err != nil {
		return err
	}

	transformation, err := iiifimage.NewTransformation(level, *region, *size, *rotation, *quality, *format)

	if err != nil {
		return err
	}

	// TO DO : compare extension of infile to 'format'

	if !transformation.HasTransformation() {
		return errors.New("No transformation")
	}

	body, err := ioutil.ReadFile(infile)

	if err != nil {
		return err
	}

	source, err := iiifsource.NewMemorySource(body)

	if err != nil {
		return err
	}

	image, err := driver.NewImageFromConfigWithSource(config, source, fname)

	if err != nil {
		return err
	}

	err = image.Transform(transformation)

	if err != nil {
		return err
	}

	fh, err := os.Create(outfile)

	if err != nil {
		return err
	}

	defer fh.Close()
	fh.Write(image.Body())
	fh.Sync()

	return nil
}
