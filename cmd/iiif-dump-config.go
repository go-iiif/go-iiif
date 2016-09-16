package main

// because this: https://github.com/thisisaaronland/go-iiif/issues/12

import (
	"flag"
	"fmt"
	iiifcompliance "github.com/thisisaaronland/go-iiif/compliance"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	"log"
)

func main() {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "example.com")

	if err != nil {
		log.Fatal(err)
	}

	compliance := level.Compliance()
	spec := compliance.Spec()

	image := spec.Image

	params := map[string]map[string]iiifcompliance.ComplianceDetails{
		"region":   image.Region,
		"size":     image.Size,
		"rotation": image.Rotation,
		"quality":  image.Quality,
		"format":   image.Format,
	}

	for p, rules := range params {

		fmt.Printf("\n### [%s](http://iiif.io/api/image/2.1/index.html#%s)\n", p, p)
		fmt.Println("| feature | syntax | required | supported |")
		fmt.Println("|---|---|---|---|")

		for feature, details := range rules {

			fmt.Printf("| %s | %s | %t | %t |\n", feature, details.Syntax, details.Required, details.Supported)
		}

	}
}
