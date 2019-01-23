package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-mimetypes"
	"log"
	"os"
	"strings"
)

func main() {

	var extension = flag.Bool("extension", false, "Lookup mimetypes by extension")
	var mimetype = flag.Bool("mimetype", false, "Lookup extensions by mimetype")

	flag.Parse()

	for _, input := range flag.Args() {

		if *mimetype {
			t := mimetypes.TypesByExtension(input)
			fmt.Printf("%s\t%s\n", input, strings.Join(t, "\t"))
		} else if *extension {
			e := mimetypes.ExtensionsByType(input)
			fmt.Printf("%s\t%s\n", input, strings.Join(e, "\t"))
		} else {
			log.Fatal("Invalid lookup type")
		}
	}

	os.Exit(0)
}
