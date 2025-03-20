package seed

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Config       *iiifconfig.Config `json:"config"`
	Mode         string             `json:"mode"`
	Endpoint     string             `json:"endpoint"`
	Quality      string             `json:"quality"`
	Format       string             `json:"format"`
	Paths        []string           `json:"paths"`
	ScaleFactors []int              `json:"scale_factors"`
	Refresh      bool               `json:"refresh"`
	Workers      int                `json:"workers"`
	NoExtension  bool               `json:"no_extension"`
	Verbose      bool               `json:"verbose"`
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign tileseed tool flags from environment variables, %w", err)
	}

	scales := make([]int, 0)

	for _, s := range strings.Split(scale_factors, ",") {

		s = strings.Trim(s, " ")
		scale, err := strconv.Atoi(s)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse scale factor, %w", err)
		}

		scales = append(scales, scale)
	}

	/*
		config, err := iiifconfig.LoadConfigWithFlagSet(ctx, fs)

		if err != nil {
			return err
		}
	*/

	paths := fs.Args()

	opts := &RunOptions{
		Mode:         mode,
		Endpoint:     endpoint,
		Quality:      quality,
		Format:       format,
		ScaleFactors: scales,
		Paths:        paths,
		Verbose:      verbose,
	}

	return opts, nil
}
