package main

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/go-iiif/go-iiif/v8/native"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"

	"github.com/go-iiif/go-iiif/v8/app/seed"
)

func main() {

	ctx := context.Background()
	err := seed.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
