package seed

import (
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Endpoint string `json:"endpoint"`
	Quality string `json:"quality"`
	Format string `json:"format"`
	Paths []string `json:"paths"`
	ScaleFactors []int `json:"scale_factors"`
	Refresh bool `json:"refresh"`
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err = flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return fmt.Errorf("Failed to assign tileseed tool flags from environment variables, %w", err)
	}

	/*

	scales := make([]int, 0)

	for _, s := range strings.Split(scale_factors, ",") {

		s = strings.Trim(s, " ")
		scale, err := strconv.Atoi(s)

		if err != nil {
			return fmt.Errorf("Failed to parse scale factor, %w", err)
		}

		scales = append(scales, scale)
	}

	*/
	
	opts := &RunOptions{}

	return opts, nil
}
