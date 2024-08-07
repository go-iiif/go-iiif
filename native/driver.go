package native

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"log/slog"

	"github.com/aaronland/go-image/colour"
	"github.com/aaronland/go-image/rotate"
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifimage "github.com/go-iiif/go-iiif/v6/image"
	iiifsource "github.com/go-iiif/go-iiif/v6/source"
	"github.com/rwcarlsen/goexif/exif"
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

	slog.Debug("NewImageFromConfigWithSource", "id", id, "source", src)

	body, err := src.Read(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read body for '%s', %w", id, err)
	}

	buf := bytes.NewBuffer(body)

	img, img_fmt, err := image.Decode(buf)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode image, %w", err)
	}

	if img_fmt == "jpeg" {

		ctx := context.Background()
		br := bytes.NewReader(body)

		o, err := rotate.GetImageOrientation(ctx, br)

		if err != nil && !exif.IsCriticalError(err) {
			return nil, fmt.Errorf("Failed to derive image orientation for '%s', %w", id, err)
		}

		if o != "0" {

			new_img, err := rotate.RotateImageWithOrientation(ctx, img, o)

			if err != nil {
				return nil, fmt.Errorf("Failed to rotate image with orientation '%s' for '%s', %w", o, id, err)
			}

			img = new_img
		}

	}

	br := bytes.NewReader(body)

	model, err := colour.DeriveModel(br)

	if err != nil {
		slog.Debug("Unable to derive model for image, default to unknown", "id", id, "error", err)
		model = colour.UnknownModel
	}

	// slog.Debug("Color model", "id", id, "mode", model)

	im := NativeImage{
		config:    config,
		source:    src,
		source_id: id,
		id:        id,
		img:       img,
		format:    img_fmt,
		model:     model,
	}

	return &im, nil
}

func (dr *NativeDriver) NewImageFromConfigWithCache(config *iiifconfig.Config, cache iiifcache.Cache, id string) (iiifimage.Image, error) {

	var image iiifimage.Image

	body, err := cache.Get(id)

	if err == nil {

		slog.Info("GOT BODY FROM CACHE", "id", id)
		source, err := iiifsource.NewMemorySource(body)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive memory source for '%s', %w", id, err)
		}

		image, err = dr.NewImageFromConfigWithSource(config, source, id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive image from source for '%s', %w", id, err)
		}

		// slog.Debug("WTF", "id", id, "model", image.ColourModel())

	} else {

		image, err = dr.NewImageFromConfig(config, id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive image from config for '%s', %w", id, err)
		}

		// slog.Debug("OMG", "id", id, "model", image.ColourModel())

		// THIS IS THE PROBLEM. WHY ARE WE ONLY CACHING image.Body which is []byte
		// and not iiifimage.Image...

		go func() {
			slog.Debug("Cache image source", "id", id)
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
