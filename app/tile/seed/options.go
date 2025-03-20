package seed

import (
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err = flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return fmt.Errorf("Failed to assign tileseed tool flags from environment variables, %w", err)
	}

	opts := &RunOptions{}

	return opts, nil
}
