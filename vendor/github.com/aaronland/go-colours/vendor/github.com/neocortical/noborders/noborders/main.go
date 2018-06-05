package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"

	"os"

	"github.com/neocortical/noborders"
)

type args struct {
	showHelp    bool
	infile      string
	outfile     string
	outfileType string
	opts        noborders.Options
}

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Printf("invalid arguments: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) <= 1 {
		fmt.Println("you must specify a file")
		os.Exit(1)
	}

	f, err := os.Open(args.infile)
	if err != nil {
		fmt.Printf("error opening a file: %s\n", err)
		os.Exit(1)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("error parsing image: %s\n", err)
		os.Exit(1)
	}

	img, err = noborders.RemoveBorders(img, args.opts)
	if err != nil {
		fmt.Printf("error processing image: %s\n", err)
		os.Exit(1)
	}
	if img == nil {
		fmt.Printf("error processing image: output image is nil\n")
		os.Exit(1)
	}

	outputFile, err := os.Create(args.outfile)
	if err != nil {
		fmt.Printf("error creating output image: %s\n", err)
		os.Exit(1)
	}

	switch args.outfileType {
	case "jpg":
		jpeg.Encode(outputFile, img, nil)
	case ".gif":
		gif.Encode(outputFile, img, nil)
	case "png":
		png.Encode(outputFile, img)
	}

	outputFile.Close()
	os.Exit(0)
}

func parseArgs() (result args, err error) {
	entropy := flag.Float64("entropy", noborders.DefaultEntropyThreshold, "Set the entropy threshold.")
	variance := flag.Float64("variance", noborders.DefaultVarianceThreshold, "Set the variance threshold.")
	multipass := flag.Bool("multipass", false, "Process the image multiple times.")
	flag.Parse()

	result.opts = noborders.Opts().
		SetEntropy(*entropy).
		SetVariance(*variance).
		SetMultiPass(*multipass)

	filenames := flag.Args()
	if len(filenames) != 2 {
		return result, errors.New("you must specify both an input and an output file")
	}

	result.infile = filenames[0]
	result.outfile = filenames[1]

	var ext = filepath.Ext(result.outfile)
	switch ext {
	case ".jpg", ".jpeg":
		result.outfileType = "jpg"
	case ".gif":
		result.outfileType = "gif"
	case ".png":
		result.outfileType = "png"
	default:
		return result, fmt.Errorf("invalid output file type: %s", ext)
	}

	return result, nil
}
