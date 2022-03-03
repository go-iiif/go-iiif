package main

// ./bin/iiif-transform -config config.json -quality dither -size ,300 /usr/local/cooperhewitt/iiif/images/184512_5f7f47e5b3c66207_x.jpg /vagrant/test2.jpg

import (
	_ "github.com/aaronland/gocloud-blob-s3"
	_ "github.com/go-iiif/go-iiif/v5/native"
	_ "gocloud.dev/blob/fileblob"
)

import (
	"context"
	"github.com/go-iiif/go-iiif/v5/tools"
	"log"
)

func main() {

	tool, err := tools.NewTransformTool()

	if err != nil {
		log.Fatal(err)
	}

	err = tool.Run(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}
