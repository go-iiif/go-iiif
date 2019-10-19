package main

import (
	_ "github.com/go-iiif/go-iiif/native"
	"github.com/go-iiif/go-iiif/tools"
	"log"
)

func main() {

	tool, err := tools.NewProcessTool()

	if err != nil {
		log.Fatal(err)
	}

	err = tool.Run()

	if err != nil {
		log.Fatal(err)
	}
}
