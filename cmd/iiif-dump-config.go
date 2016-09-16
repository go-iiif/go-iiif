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

	//

	type FeatureDetails struct {
		feature          string
		syntax           string
		required_spec    bool
		supported_spec   bool
		required_actual  bool
		supported_actual bool
	}

	fd := make(map[string]map[string]FeatureDetails)

	//

	spec, err := iiifcompliance.NewLevel2ComplianceSpec()

	if err != nil {
		log.Fatal(err)
	}

	image := spec.Image

	params := map[string]map[string]iiifcompliance.ComplianceDetails{
		"region":   image.Region,
		"size":     image.Size,
		"rotation": image.Rotation,
		"quality":  image.Quality,
		"format":   image.Format,
	}

	//

	for p, rules := range params {

		fd[p] = make(map[string]FeatureDetails)

		for feature, details := range rules {

			fd[p][feature] = FeatureDetails{
				feature:          feature,
				syntax:           details.Syntax,
				required_spec:    details.Required,
				supported_spec:   details.Supported,
				required_actual:  details.Required,
				supported_actual: details.Supported,
			}
		}
	}

	//

	compliance := level.Compliance()
	actual := compliance.Spec()

	image = actual.Image

	params = map[string]map[string]iiifcompliance.ComplianceDetails{
		"region":   image.Region,
		"size":     image.Size,
		"rotation": image.Rotation,
		"quality":  image.Quality,
		"format":   image.Format,
	}

	for p, rules := range params {

		for feature, details := range rules {

			_f := fd[p][feature]
			_f.required_actual = details.Required
			_f.supported_actual = details.Supported

			fd[p][feature] = _f
		}
	}

	//

	prms := []string{
		"region", "size", "rotation", "quality", "format",
	}

	for _, p := range prms {

		rules := fd[p]

		fmt.Printf("\n### [%s](http://iiif.io/api/image/2.1/index.html#%s)\n", p, p)
		fmt.Printf("| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |\n")
		fmt.Printf("|---|---|---|---|---|---|\n")

		for feature, details := range rules {

			fmt.Printf("| **%s** | %s | %t | %t | **%t** | **%t** |\n", feature, details.syntax, details.required_spec, details.supported_spec, details.required_actual, details.supported_actual)
		}

	}

}
