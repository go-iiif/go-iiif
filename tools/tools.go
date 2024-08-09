package tools

// Important: Almost everything in the package will change in the v7 release.
// Not so much the functionality the way in which the code is organized. "Tools"
// will no longer be concered with flags - those will be moved in to the app
// package - but instead only use tool-specific "options" structs. The Tool interface
// will probably only contain a single Run(ctx context.Context, uris ...string) error
// method.

import (
	"context"
	"flag"
)

type Tool interface {
	Run(context.Context) error
	RunWithFlagSet(context.Context, *flag.FlagSet) error
	RunWithFlagSetAndPaths(context.Context, *flag.FlagSet, ...string) error
}
