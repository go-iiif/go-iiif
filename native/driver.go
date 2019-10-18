package native

import (
	"bytes"
	_ "fmt"
	iiif "github.com/go-iiif/go-iiif"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiifsource "github.com/go-iiif/go-iiif/source"
	_ "image"
	"image/gif"
	_ "log"
)

func init() {
	dr, err := NewNativeDriver()

	if err != nil {
		panic(err)
	}

	iiif.RegisterDriver("native", dr)
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

	img, _, err := gif.Decode(buf) // FIX ME...

	if err != nil {
		return nil, err
	}

	im := NativeImage{
		config:    config,
		source:    src,
		source_id: id,
		id:        id,
		img:       img,
		isgif:     false,
	}

	/*

		Hey look - see the 'isgif' flag? We're going to hijack the fact that
		img doesn't handle GIF files and if someone requests them then we
		will do the conversion after the final call to im.img.Process and
		after we do handle any custom features. We are relying on the fact
		that both img.NewImage and img.Image() expect and return raw bytes
		and we are ignoring whatever img thinks in the Format() function.
		So basically you should not try to any processing in img/libvips
		after the -> GIF transformation. (20160922/thisisaaronland)

		See also: https://github.com/h2non/img/issues/41
	*/

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

		image, err = NewImageFromConfigWithSource(config, source, id)

		if err != nil {
			return nil, err
		}

	} else {

		image, err = NewImageFromConfig(config, id)

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
