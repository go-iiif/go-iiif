package main

import (
	"flag"
	"fmt"
	"github.com/go-iiif/go-iiif-uri"
	"log"
)

func main() {

	var uri_type = flag.String("type", "string", "...")

	flag.Parse()

	for _, str_uri := range flag.Args() {

		u, err := uri.NewURIWithType(str_uri, *uri_type)

		if err != nil {
			msg := fmt.Sprintf("Invalid URI (%s) %s", str_uri, err)
			log.Fatal(msg)
		}

		log.Printf("%s OK\n", u.String())
	}
}
