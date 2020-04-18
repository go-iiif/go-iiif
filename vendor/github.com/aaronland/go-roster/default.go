package roster

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
)

type DefaultRoster struct {
	Roster
	mu      *sync.RWMutex
	drivers map[string]interface{}
}

func NewDefaultRoster() (Roster, error) {

	mu := new(sync.RWMutex)
	drivers := make(map[string]interface{})

	dr := &DefaultRoster{
		mu:      mu,
		drivers: drivers,
	}

	return dr, nil
}

func (dr *DefaultRoster) Driver(ctx context.Context, name string) (interface{}, error) {

	nrml_name := dr.NormalizeName(ctx, name)

	dr.mu.Lock()
	defer dr.mu.Unlock()

	i, ok := dr.drivers[nrml_name]

	if !ok {
		return nil, errors.New("Unknown driver")
	}

	return i, nil
}

func (dr *DefaultRoster) Register(ctx context.Context, name string, i interface{}) error {

	dr.mu.Lock()
	defer dr.mu.Unlock()

	if i == nil {
		return errors.New("Nothing to register")
	}

	nrml_name := dr.NormalizeName(ctx, name)

	_, dup := dr.drivers[nrml_name]

	if dup {
		return errors.New("Register called twice for reader " + name)
	}

	dr.drivers[nrml_name] = i
	return nil
}

func (dr *DefaultRoster) UnregisterAll(ctx context.Context) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	dr.drivers = make(map[string]interface{})
	return nil
}

func (dr *DefaultRoster) NormalizeName(ctx context.Context, name string) string {
	return strings.ToUpper(name)
}

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
