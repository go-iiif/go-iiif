package main

import (
	"github.com/koyachi/go-atkinson"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

func main() {
	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	img, err := atkinson.DitherFile(path)
	if err != nil {
		log.Fatal(err)
	}

	path, err = filepath.Abs(path + ".result.jpg")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(file, img, nil)
	if err != nil {
		log.Fatal(err)
	}
}
