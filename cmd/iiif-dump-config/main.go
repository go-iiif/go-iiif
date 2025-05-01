package main

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/go-iiif/go-iiif/v8/native"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"

	"github.com/go-iiif/go-iiif/v8/app/config/dump"
)

func main() {

	ctx := context.Background()
	err := dump.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
