package native

import (
	"bytes"
	_ "fmt"
	iiifcache "github.com/go-iiif/go-iiif/v2/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v2/config"
	iiifdriver "github.com/go-iiif/go-iiif/v2/driver"
	iiifimage "github.com/go-iiif/go-iiif/v2/image"
	iiifsource "github.com/go-iiif/go-iiif/v2/source"
	"image"
	_ "log"
)

func init() {

	dr, err := NewNativeDriver()

	if err != nil {
		panic(err)
	}

	iiifdriver.RegisterDriver("native", dr)
}

type NativeDriver struct {
	iiifdriver.Driver
}

func NewNativeDriver() (iiifdriver.Driver, error) {
	dr := &NativeDriver{}
	return dr, nil
}

func (dr *NativeDriver) NewImageFromConfigWithSource(config *iiifconfig.Config, src iiifsource.Source, id string) (iiifimage.Image, error) {

	body, err := src.Read(id)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(body)

	img, fmt, err := image.Decode(buf)

	if err != nil {
		return nil, err
	}

	im := NativeImage{
		config:    config,
		source:    src,
		source_id: id,
		id:        id,
		img:       img,
		format:    fmt,
	}

	return &im, nil
}

func (dr *NativeDriver) NewImageFromConfigWithCache(config *iiifconfig.Config, cache iiifcache.Cache, id string) (iiifimage.Image, error) {

	var image iiifimage.Image

	body, err := cache.Get(id)

	if err == nil {

		source, err := iiifsource.NewMemorySource(body)

		if err != nil {
			return nil, err
		}

		image, err = dr.NewImageFromConfigWithSource(config, source, id)

		if err != nil {
			return nil, err
		}

	} else {

		image, err = dr.NewImageFromConfig(config, id)

		if err != nil {
			return nil, err
		}

		go func() {
			cache.Set(id, image.Body())
		}()
	}

	return image, nil
}

func (dr *NativeDriver) NewImageFromConfig(config *iiifconfig.Config, id string) (iiifimage.Image, error) {

	source, err := iiifsource.NewSourceFromConfig(config)

	if err != nil {
		return nil, err
	}

	return dr.NewImageFromConfigWithSource(config, source, id)
}
