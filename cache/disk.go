package cache

import (
	"github.com/thisisaaronland/go-iiif/config"
	"io/ioutil"
	_ "log"
	"os"
	"path"
	"path/filepath"
)

type DiskCache struct {
	Cache
	root string
}

func NewDiskCache(cfg config.CacheConfig) (*DiskCache, error) {

	root := cfg.Path
	_, err := os.Stat(root)

	if os.IsNotExist(err) {
		return nil, err
	}

	c := DiskCache{
		root: root,
	}

	return &c, nil
}

func (c *DiskCache) Exists(rel_path string) bool {

	abs_path := path.Join(c.root, rel_path)

	_, err := os.Stat(abs_path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (c *DiskCache) Get(rel_path string) ([]byte, error) {

	abs_path := path.Join(c.root, rel_path)

	_, err := os.Stat(abs_path)

	if os.IsNotExist(err) {
		// fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadFile(abs_path)

	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	return body, nil
}

func (c *DiskCache) Set(rel_path string, body []byte) error {

	abs_path := path.Join(c.root, rel_path)

	root := filepath.Dir(abs_path)

	_, err := os.Stat(root)

	if os.IsNotExist(err) {
		os.MkdirAll(root, 0755)
	}

	fh, err := os.Create(abs_path)

	if err != nil {
		return err
	}

	defer fh.Close()
	fh.Write(body)
	fh.Sync()

	return nil
}

func (c *DiskCache) Unset(rel_path string) error {

	abs_path := path.Join(c.root, rel_path)

	_, err := os.Stat(abs_path)

	if os.IsNotExist(err) {
		return nil
	}

	return os.Remove(abs_path)
}
