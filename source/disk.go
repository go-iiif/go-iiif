package source

import (
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DiskSource struct {
	Source
	root string
}

func NewDiskSource(config *iiifconfig.Config) (*DiskSource, error) {

	cfg := config.Images

	ds := DiskSource{
		root: cfg.Source.Path,
	}

	return &ds, nil
}

func (ds *DiskSource) Read(uri string) ([]byte, error) {

	abs_path := filepath.Join(ds.root, uri)

	_, err := os.Stat(abs_path)

	if os.IsNotExist(err) {
		return nil, err
	}

	body, err := ioutil.ReadFile(abs_path)

	if err != nil {
		return nil, err
	}

	return body, nil
}
