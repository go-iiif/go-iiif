package main

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/go-iiif/go-iiif/v7/native"
	"github.com/go-iiif/go-iiif/v7/tools"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
)

func main() {

	tool, err := tools.NewTileSeedTool()

	if err != nil {
		log.Fatal(err)
	}

	err = tool.Run(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}
