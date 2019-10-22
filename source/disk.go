package source

import (
	"fmt"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"
)

type DiskSource struct {
	Source
	root string
}

func NewDiskSource(config *iiifconfig.Config) (Source, error) {

	cfg := config.Images
	uri := fmt.Sprintf("file://%s", cfg.Source.Path)

	return NewBlobSourceFromURI(uri)

	// PLEASE REMOVE EVERYTHING ELSE AS SOON AS POSSIBLE

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
