package main

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/go-iiif/go-iiif/v6/native"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"

	"github.com/go-iiif/go-iiif/v6/app/tile/seed"
)

func main() {

	ctx := context.Background()
	err := seed.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
