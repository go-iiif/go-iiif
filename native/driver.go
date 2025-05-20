package native

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/aaronland/go-image/v2/colour"
	"github.com/aaronland/go-image/v2/decode"
	"github.com/aaronland/go-image/v2/rotate"
	iiifcache "github.com/go-iiif/go-iiif/v8/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiifdriver "github.com/go-iiif/go-iiif/v8/driver"
	iiifimage "github.com/go-iiif/go-iiif/v8/image"
	iiifsource "github.com/go-iiif/go-iiif/v8/source"
	"github.com/rwcarlsen/goexif/exif"
)

func init() {

	err := iiifdriver.RegisterDriver(context.Background(), "native", NewNativeDriver)

	if err != nil {
		panic(err)
	}
}

type NativeDriver struct {
	iiifdriver.Driver
}

func NewNativeDriver(ctx context.Context, uri string) (iiifdriver.Driver, error) {
	dr := &NativeDriver{}
	return dr, nil
}

func (dr *NativeDriver) NewImageFromConfigWithSource(ctx context.Context, config *iiifconfig.Config, src iiifsource.Source, id string) (iiifimage.Image, error) {

	logger := slog.Default()
	logger = logger.With("source", src)
	logger = logger.With("id", id)

	// logger.Debug("New image from config with source")

	body, err := src.Read(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read body for '%s', %w", id, err)
	}

	buf := bytes.NewReader(body)

	img, img_fmt, _, err := decode.DecodeImage(ctx, buf)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode image, %w", err)
	}

	_, err = buf.Seek(0, 0)

	if err != nil {
		return nil, fmt.Errorf("Failed to rewind buffer, %w", err)
	}

	model, err := colour.DeriveModel(buf)

	if err != nil {
		slog.Debug("Unable to derive model for image, default to unknown", "id", id, "error", err)
		model = colour.UnknownModel
	}

	// logger.Debug("New Golang image", "format", img_fmt, "model", model)

	switch model {
	case colour.AppleDisplayP3Model:
		img = colour.ToDisplayP3(img)
	case colour.AdobeRGBModel:
		img = colour.ToAdobeRGB(img)
	case colour.UnknownModel, colour.SRGBModel:
		// pass
	default:
		// pass
	}

	if img_fmt == "jpeg" {

		_, err = buf.Seek(0, 0)

		if err != nil {
			return nil, fmt.Errorf("Failed to rewind buffer, %w", err)
		}

		ctx := context.Background()

		o, err := rotate.GetImageOrientation(ctx, buf)

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

func (dr *NativeDriver) NewImageFromConfigWithCache(ctx context.Context, config *iiifconfig.Config, cache iiifcache.Cache, id string) (iiifimage.Image, error) {

	var image iiifimage.Image

	body, err := cache.Get(id)

	if err == nil {

		source, err := iiifsource.NewMemorySourceWithKey(id, body)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive memory source for '%s', %w", id, err)
		}

		image, err = dr.NewImageFromConfigWithSource(ctx, config, source, id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive image from source for '%s', %w", id, err)
		}

	} else {

		image, err = dr.NewImageFromConfig(ctx, config, id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive image from config for '%s', %w", id, err)
		}

		go func() {
			slog.Debug("Cache image source", "id", id)
			cache.Set(id, image.Body())
		}()
	}

	return image, nil
}

func (dr *NativeDriver) NewImageFromConfig(ctx context.Context, cfg *iiifconfig.Config, id string) (iiifimage.Image, error) {

	source, err := iiifsource.NewSource(ctx, cfg.Images.Source.URI)

	if err != nil {
		return nil, err
	}

	return dr.NewImageFromConfigWithSource(ctx, cfg, source, id)
}
