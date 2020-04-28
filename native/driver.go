package native

import (
	"bytes"
	"context"
	_ "fmt"
	"github.com/aaronland/go-image-rotate"
	iiifcache "github.com/go-iiif/go-iiif/v4/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiifdriver "github.com/go-iiif/go-iiif/v4/driver"
	iiifimage "github.com/go-iiif/go-iiif/v4/image"
	iiifsource "github.com/go-iiif/go-iiif/v4/source"	
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

	// ROTATION STUFF GOES HERE
	// "github.com/aaronland/go-image-rotate"

	if fmt == "jpeg" {

		ctx := context.Background()		
		br := bytes.NewReader(body)
		
		o, err := rotate.GetImageOrientation(ctx, br)

		if err != nil {
			return nil, err
		}
		
		new_img, err := rotate.RotateImageWithOrientation(ctx, img, o)

		if err != nil {
			return nil, err
		}

		img = new_img
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
