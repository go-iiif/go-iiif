package roster

import (
	"context"
)

// type Roster is an interface for defining internal lookup tables (or "rosters") for registering and instantiatinge custom interfaces with multiple implementations.
type Roster interface {
	// Driver returns the value associated with a name or scheme from the list of drivers that have been registered.
	Driver(context.Context, string) (interface{}, error)
	// Drivers returns the list of names or schemes for the list of drivers that have been registered.
	Drivers(context.Context) []string
	// UnregisterAll removes all the registered drivers from an instance implementing the Roster interfave.
	UnregisterAll(context.Context) error
	// NormalizeName returns a normalized version of a string.
	NormalizeName(context.Context, string) string
	// Register associated a name or scheme with an arbitrary interface.
	Register(context.Context, string, interface{}) error
}
