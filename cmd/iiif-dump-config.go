package main

// because this: https://github.com/thisisaaronland/go-iiif/issues/12

import (
	"flag"
	"fmt"
	iiifcompliance "github.com/thisisaaronland/go-iiif/compliance"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	"log"
	"strings"
)

func Format(spec *iiifcompliance.Level2ComplianceSpec) string {

	rsp := ""

	image := spec.Image

	params := map[string]map[string]iiifcompliance.ComplianceDetails{
		"region":   image.Region,
		"size":     image.Size,
		"rotation": image.Rotation,
		"quality":  image.Quality,
		"format":   image.Format,
	}

	for p, rules := range params {

		rsp += fmt.Sprintf("\n### [%s](http://iiif.io/api/image/2.1/index.html#%s)\n", p, p)
		rsp += fmt.Sprintf("| feature | syntax | required | supported |\n")
		rsp += fmt.Sprintf("|---|---|---|---|\n")

		for feature, details := range rules {

			rsp += fmt.Sprintf("| %s | %s | %t | %t |\n", feature, details.Syntax, details.Required, details.Supported)
		}

	}

	return rsp
}

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
	current := compliance.Spec()

	ideal, err := iiifcompliance.NewLevel2ComplianceSpec()

	if err != nil {
		log.Fatal(err)
	}

	configs := map[string]*iiifcompliance.Level2ComplianceSpec{
		"default": ideal,
		"current": current,
	}

	labels := make([]string, 0)
	details := make([]string, 0)

	for name, cfg := range configs {
		labels = append(labels, name)
		details = append(details, Format(cfg))
	}

	str_labels := strings.Join(labels, " | ")
	str_details := strings.Join(details, " | ")

	fmt.Printf("| %s |\n", str_labels)
	fmt.Printf("|---|---|\n")
	fmt.Printf("| %s |\n", str_details)
}
