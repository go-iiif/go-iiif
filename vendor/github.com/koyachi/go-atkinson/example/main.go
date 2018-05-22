package main

import (
	"../"
	"fmt"
	"github.com/koyachi/go-lena"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fmt.Printf("processing...\n")
	img, err := lena.Image()
	if err != nil {
		log.Fatal(err)
	}
	img, err = atkinson.Dither(img)
	if err != nil {
		log.Fatal(err)
	}

	path, err := filepath.Abs("result.jpg")
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
