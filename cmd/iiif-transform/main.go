package main

// ./bin/iiif-transform -config config.json -quality dither -size ,300 /usr/local/cooperhewitt/iiif/images/184512_5f7f47e5b3c66207_x.jpg /vagrant/test2.jpg

import (
	_ "github.com/go-iiif/go-iiif/native"
	"github.com/go-iiif/go-iiif/tools"
	"log"
)

func main() {

	tool, err := tools.NewTransformTool()

	if err != nil {
		log.Fatal(err)
	}

	err = tool.Run()

	if err != nil {
		log.Fatal(err)
	}

}
