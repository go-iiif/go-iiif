package uri

import (
	"errors"
	_ "log"
	"net/url"
	"sort"
	"strings"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

type Driver interface {
	NewURI(string) (URI, error)
}

func NewURIWithDriver(str_uri string) (URI, error) {

	driversMu.Lock()
	defer driversMu.Unlock()

	u, err := url.Parse(str_uri)

	if err != nil {
		return nil, err
	}

	name := u.Scheme
	name_nrml := normalizeDriverName(name)

	driver, ok := drivers[name_nrml]

	if !ok {
		return nil, errors.New("Unknown driver")
	}

	return driver.NewURI(str_uri)
}

func RegisterDriver(name string, driver Driver) {

	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("go-iiif-uri: Register driver is nil")

	}

	name_nrml := normalizeDriverName(name)

	if _, dup := drivers[name_nrml]; dup {
		panic("index: Register called twice for driver " + name_nrml)
	}

	drivers[name_nrml] = driver
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

func normalizeDriverName(name string) string {
	return strings.ToLower(name)
}
