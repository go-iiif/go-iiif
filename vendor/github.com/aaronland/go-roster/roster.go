package roster

import (
	"context"
)

type Roster interface {
	Driver(context.Context, string) (interface{}, error)
	Drivers(context.Context) []string
	UnregisterAll(context.Context) error
	NormalizeName(context.Context, string) string
	Register(context.Context, string, interface{}) error
}
