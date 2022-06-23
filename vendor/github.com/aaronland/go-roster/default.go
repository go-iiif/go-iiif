package roster

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// DefaultRoster implements the the `Roster` interface mapping scheme names to arbitrary interface values.
type DefaultRoster struct {
	Roster
	mu      *sync.RWMutex
	drivers map[string]interface{}
}

// NewDefaultRoster returns a new `DefaultRoster` instance.
func NewDefaultRoster() (Roster, error) {

	mu := new(sync.RWMutex)
	drivers := make(map[string]interface{})

	dr := &DefaultRoster{
		mu:      mu,
		drivers: drivers,
	}

	return dr, nil
}

// Driver returns the value associated with the key for the normalized value of 'name' in the list of registered
// drivers available to 'dr'.
func (dr *DefaultRoster) Driver(ctx context.Context, name string) (interface{}, error) {

	nrml_name := dr.NormalizeName(ctx, name)

	dr.mu.Lock()
	defer dr.mu.Unlock()

	i, ok := dr.drivers[nrml_name]

	if !ok {
		return nil, fmt.Errorf("Unknown driver: %s (%s)", name, nrml_name)
	}

	return i, nil
}

// Registers creates a new entry in the list of drivers available to 'dr' mapping the normalized version of 'name' to 'i'.
func (dr *DefaultRoster) Register(ctx context.Context, name string, i interface{}) error {

	dr.mu.Lock()
	defer dr.mu.Unlock()

	if i == nil {
		return errors.New("Nothing to register")
	}

	nrml_name := dr.NormalizeName(ctx, name)

	_, dup := dr.drivers[nrml_name]

	if dup {
		return fmt.Errorf("Register called twice for reader '%s'", name)
	}

	dr.drivers[nrml_name] = i
	return nil
}

// UnregisterAll removes all the registers drivers from 'dr'.
func (dr *DefaultRoster) UnregisterAll(ctx context.Context) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	dr.drivers = make(map[string]interface{})
	return nil
}

// NormalizeName returns a normalized (upper-cased) version of 'name'.
func (dr *DefaultRoster) NormalizeName(ctx context.Context, name string) string {
	return strings.ToUpper(name)
}

// Drivers returns the list of registered schemes for all the drivers available to 'dr'.
func (dr *DefaultRoster) Drivers(ctx context.Context) []string {

	dr.mu.RLock()
	defer dr.mu.RUnlock()

	var list []string

	for name := range dr.drivers {
		list = append(list, name)
	}

	sort.Strings(list)
	return list
}
