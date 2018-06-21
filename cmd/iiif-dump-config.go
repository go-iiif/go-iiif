package main

// because this: https://github.com/thisisaaronland/go-iiif/issues/12

import (
	"flag"
	"fmt"
	iiifcompliance "github.com/thisisaaronland/go-iiif/compliance"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	"log"
	"sort"
)

type FeatureDetails struct {
	feature          string
	syntax           string
	required_spec    bool
	supported_spec   bool
	required_config  bool
	supported_config bool
}

func Sorted(h map[string]FeatureDetails) []string {

	keys := make([]string, 0)

	for k, _ := range h {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func main() {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "example.com")

	if err != nil {
		log.Fatal(err)
	}

	//

	fd := make(map[string]map[string]FeatureDetails)

	//

	spec, err := iiifcompliance.NewLevel2ComplianceSpec()

	if err != nil {
		log.Fatal(err)
	}

	//

	var image iiifcompliance.ImageCompliance // because we're going to instantiate this twice with two different values

	compliance := level.Compliance()
	actual := compliance.Spec()

	image = actual.Image

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
				required_config:  details.Required,
				supported_config: details.Supported,
			}
		}
	}

	//

	image = spec.Image

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
			_f.required_spec = details.Required
			_f.supported_spec = details.Supported

			fd[p][feature] = _f
		}
	}

	//

	prms := []string{
		"region", "size", "rotation", "quality", "format",
	}

	for _, p := range prms {

		rules := fd[p]

		fmt.Printf("\n##### [%s](http://iiif.io/api/image/2.1/index.html#%s)\n", p, p)
		fmt.Printf("| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |\n")
		fmt.Printf("|---|---|---|---|---|---|\n")

		features := Sorted(rules)

		for _, feature := range features {

			details := rules[feature]

			rs := "green"
			ss := "green"
			rc := "green"
			sc := "green"

			if !details.required_spec {
				rs = "red"
			}

			if !details.supported_spec {
				ss = "red"
			}

			if !details.required_config {
				rc = "red"
			}

			if !details.supported_config {
				sc = "red"
			}

			rs_html := fmt.Sprintf("<span style=\"color:%s;\">%t</span>", rs, details.required_spec)
			ss_html := fmt.Sprintf("<span style=\"color:%s;\">%t</span>", ss, details.supported_spec)

			rc_html := fmt.Sprintf("<span style=\"color:%s;\">%t</span>", rc, details.required_config)
			sc_html := fmt.Sprintf("<span style=\"color:%s;\">%t</span>", sc, details.supported_config)

			if details.required_config {
				rc_html = fmt.Sprintf("<span style=\"color:%s;\">**%t**</span>", rc, details.required_config)
			}

			if details.supported_config {
				sc_html = fmt.Sprintf("<span style=\"color:%s;\">**%t**</span>", sc, details.supported_config)
			}

			fmt.Printf("| **%s** | %s | %s | %s | %s | %s |\n", feature, details.syntax, rs_html, ss_html, rc_html, sc_html)
		}

	}

}
