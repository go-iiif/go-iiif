package main

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/go-iiif/go-iiif/v6/native"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"

	"github.com/go-iiif/go-iiif/v6/app/server"
)

func main() {

	ctx := context.Background()
	err := server.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
