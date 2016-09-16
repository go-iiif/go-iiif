package main

// because this: https://github.com/thisisaaronland/go-iiif/issues/12

import (
	"flag"
	"fmt"
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

	fmt.Println("\n### region\n")
	fmt.Println("| feature | syntax | required | supported |")
	fmt.Println("|---|---|---|---|")

	for feature, details := range image.Region {

		fmt.Printf("| %s | %s | %t | %t |\n", feature, details.Syntax, details.Required, details.Supported)
	}

	fmt.Println("\n### size\n")
	fmt.Println("| feature | syntax | required | supported |")
	fmt.Println("|---|---|---|---|")

	for feature, details := range image.Size {

		fmt.Printf("| %s | %s | %t | %t |\n", feature, details.Syntax, details.Required, details.Supported)
	}

	fmt.Println("\n### rotation\n")
	fmt.Println("| feature | syntax | required | supported |")
	fmt.Println("|---|---|---|---|")

	for feature, details := range image.Rotation {

		fmt.Printf("| %s | %s | %t | %t |\n", feature, details.Syntax, details.Required, details.Supported)
	}

	fmt.Println("\n### quality\n")
	fmt.Println("| feature | syntax | required | supported |")
	fmt.Println("|---|---|---|---|")

	for feature, details := range image.Quality {

		fmt.Printf("| %s | %s | %t | %t |\n", feature, details.Syntax, details.Required, details.Supported)
	}

	fmt.Println("\n### format\n")
	fmt.Println("| feature | syntax | required | supported |")
	fmt.Println("|---|---|---|---|")

	for feature, details := range image.Format {

		fmt.Printf("| %s | %s | %t | %t |\n", feature, details.Syntax, details.Required, details.Supported)
	}

}
