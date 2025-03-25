package main

// ./bin/iiif-transform -config config.json -quality dither -size ,300 /usr/local/cooperhewitt/iiif/images/184512_5f7f47e5b3c66207_x.jpg /vagrant/test2.jpg

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/go-iiif/go-iiif/v6/native"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"	

	"github.com/go-iiif/go-iiif/v6/app/transform"
)

func main() {

	ctx := context.Background()
	err := transform.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
