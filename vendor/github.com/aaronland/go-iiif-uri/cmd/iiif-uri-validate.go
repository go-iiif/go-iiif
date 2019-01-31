package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/go-iiif-uri"
	"log"
)

func main() {

	var uri_type = flag.String("type", "iiif", "...")
	
	flag.Parse()

	for _, str_uri := range flag.Args() {

		var u uri.URI
		var e error
		
		switch *uri_type {

		case "iiif":
			u, e = uri.NewIIIFURI(str_uri)
			case "rewrite":
			u, e = uri.NewRewriteURI(str_uri)
		default:
			e = errors.New("Unknown URI type")
		}

		if e != nil {
			msg := fmt.Sprintf("Invalid URI (%s) %s", str_uri, e)
			log.Fatal(msg)
		}

		log.Printf("%s OK\n", u.String())
	}
}
