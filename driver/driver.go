package driver

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifimage "github.com/go-iiif/go-iiif/v6/image"
	iiifsource "github.com/go-iiif/go-iiif/v6/source"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

type Driver interface {
	NewImageFromConfigWithSource(*iiifconfig.Config, iiifsource.Source, string) (iiifimage.Image, error)
	NewImageFromConfigWithCache(*iiifconfig.Config, iiifcache.Cache, string) (iiifimage.Image, error)
	NewImageFromConfig(*iiifconfig.Config, string) (iiifimage.Image, error)
}

func RegisterDriver(name string, driver Driver) {

	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("iiif: Register driver is nil")

	}

	nrml_name := normalizeName(name)

	if _, dup := drivers[nrml_name]; dup {
		panic("index: Register called twice for driver " + name)
	}

	drivers[nrml_name] = driver
}

func normalizeName(name string) string {
	return strings.ToUpper(name)
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	drivers = make(map[string]Driver)
}

func Drivers() []string {

	driversMu.RLock()
	defer driversMu.RUnlock()

	var list []string

	for name := range drivers {
		list = append(list, name)
	}

	sort.Strings(list)
	return list
}

func NewDriver(name string) (Driver, error) {

	driversMu.RLock()
	defer driversMu.RUnlock()

	nrml_name := normalizeName(name)

	dr, ok := drivers[nrml_name]

	if !ok {
		msg := fmt.Sprintf("Invalid go-iiif driver '%s' ('%s')", name, nrml_name)
		return nil, errors.New(msg)
	}

	return dr, nil
}

func NewDriverFromConfig(config *iiifconfig.Config) (Driver, error) {
	return NewDriver(config.Graphics.Source.Name)
}
