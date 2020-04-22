package main

import (
	"context"
	_ "github.com/aaronland/go-cloud-s3blob"
	_ "github.com/go-iiif/go-iiif/v3/native"
	"github.com/go-iiif/go-iiif/v3/tools"
	_ "gocloud.dev/blob/fileblob"
	"log"
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
