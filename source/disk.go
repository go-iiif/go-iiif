package source

import (
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DiskSource struct {
	Source
	root  string
	cache iiifcache.Cache
}

func NewDiskSource(config iiifconfig.Config) (*DiskSource, error) {

	cfg := config.Images

	ch, err := iiifcache.NewCacheFromConfig(config)

	if err != nil {
		return nil, err
	}

	ds := DiskSource{
		root:  cfg.Source.Path,
		cache: ch,
	}

	return &ds, nil
}

func (ds *DiskSource) Read(uri string) ([]byte, error) {

	body, err := ds.cache.Get(uri)

	if err == nil {
		return body, nil
	}

	abs_path := filepath.Join(ds.root, uri)

	_, err = os.Stat(abs_path)

	if os.IsNotExist(err) {
		return nil, err
	}

	body, err = ioutil.ReadFile(abs_path)

	if err != nil {
		return nil, err
	}

	go func() {
		ds.cache.Set(uri, body)
	}()

	return body, nil
}
