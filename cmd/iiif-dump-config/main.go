package main

// because this: https://github.com/go-iiif/go-iiif/issues/12

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"

	_ "github.com/aaronland/gocloud-blob/s3"
	iiifcompliance "github.com/go-iiif/go-iiif/v6/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiiflevel "github.com/go-iiif/go-iiif/v6/level"
	iiiftools "github.com/go-iiif/go-iiif/v6/tools"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
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

	ctx := context.Background()

	fs := flag.NewFlagSet("dump", flag.ExitOnError)
	err := iiiftools.AppendCommonConfigFlags(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to append config flags, %v", err)
	}

	flagset.Parse(fs)

	config_source, err := lookup.StringVar(fs, "config-source")

	if err != nil {
		log.Fatalf("Failed to parse -config-source flag, %v", err)
	}

	config_name, err := lookup.StringVar(fs, "config-name")

	if err != nil {
		log.Fatalf("Failed to parse -config-name flag, %v", err)
	}

	config_bucket, err := blob.OpenBucket(ctx, config_source)

	if err != nil {
		log.Fatalf("Failed to open config bucket, %v", err)
	}

	config, err := iiifconfig.NewConfigFromBucket(ctx, config_bucket, config_name)

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
