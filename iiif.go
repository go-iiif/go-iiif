package iiif

import (
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	"sort"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]iiifdriver.Driver)
)

func RegisterDriver(name string, driver iiifdriver.Driver) {

	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("iiif: Register driver is nil")

	}

	if _, dup := drivers[name]; dup {
		panic("index: Register called twice for driver " + name)
	}

	drivers[name] = driver
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	drivers = make(map[string]iiifdriver.Driver)
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
