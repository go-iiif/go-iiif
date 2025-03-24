package image

import (
	"bytes"
	"fmt"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log/slog"
	"time"

	"github.com/aaronland/go-image/colour"
	"github.com/dgraph-io/ristretto/v2"
)

var golang_image_cache *ristretto.Cache[string, image.Image]

func init() {

	cache, err := ristretto.NewCache(&ristretto.Config[string, image.Image]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		slog.Error("Failed to set up golang image cache", "error", err)
		return
	}

	golang_image_cache = cache
}

type Image interface {
	Identifier() string
	Rename(string) error
	Transform(*Transformation) error // http://iiif.io/api/image/2.1/#order-of-implementation
	Update([]byte) error
	Body() []byte
	Format() string
	ContentType() string
	Dimensions() (Dimensions, error)
	ColourModel() colour.Model
}

type Dimensions interface {
	Height() int
	Width() int
}

// Convert a go-iiif/image.Image instance to a Go language image.Image instance.
func IIIFImageToGolangImage(im Image) (image.Image, error) {

	logger := slog.Default()
	logger = logger.With("id", im.Identifier())

	t1 := time.Now()
	defer func() {
		logger.Info("Time to transform image", "time", time.Since(t1))
	}()

	if golang_image_cache != nil {

		go_im, ok := golang_image_cache.Get(im.Identifier())

		if ok {
			logger.Debug("Cache HIT")
			return go_im, nil
		}

		logger.Debug("Cache MISS")
	}

	var goimg image.Image
	var err error

	content_type := im.ContentType()

	if content_type == "image/gif" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = gif.Decode(buf)

	} else if content_type == "image/jpeg" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = jpeg.Decode(buf)

	} else if content_type == "image/png" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = png.Decode(buf)

	} else if content_type == "image/tiff" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = tiff.Decode(buf)

	} else if content_type == "image/webp" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = webp.Decode(buf)

	} else {
		err = fmt.Errorf("Unsupported content type '%s' for decoding", content_type)
	}

	if err != nil {
		logger.Error("Failed to convert image", "error", err)
		return nil, err
	}

	if golang_image_cache != nil {

		ttl := 2 * time.Minute
		ok := golang_image_cache.SetWithTTL(im.Identifier(), goimg, 1, ttl)

		if !ok {
			logger.Error("Failed to set cache for image", "error", err)
		}
	}

	return goimg, nil
}

// Assign a Go language image.Image instance to a go-iiif/image.Image instance.
func GolangImageToIIIFImage(goimg image.Image, im Image) error {

	body, err := GolangImageToBytes(goimg, im.ContentType())

	if err != nil {
		return err
	}

	return im.Update(body)
}

// Encode a Go language image.Image instance to a byte array.
func GolangImageToBytes(goimg image.Image, content_type string) ([]byte, error) {

	var out *bytes.Buffer
	var err error

	if content_type == "image/gif" {

		/*
			opts := gif.Options{
				NumColors: 256,
			}
		*/

		out = new(bytes.Buffer)
		err = gif.Encode(out, goimg, nil)

	} else if content_type == "image/jpeg" {

		out = new(bytes.Buffer)
		err = jpeg.Encode(out, goimg, nil)

	} else if content_type == "image/png" {

		out = new(bytes.Buffer)
		err = png.Encode(out, goimg)

	} else if content_type == "image/tiff" {

		out = new(bytes.Buffer)
		err = tiff.Encode(out, goimg, nil)

	} else {

		err = fmt.Errorf("Unsupported content type '%s' for encoding", content_type)
	}

	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func ApplyColourModel(im image.Image, model colour.Model) image.Image {

	switch model {
	case colour.AppleDisplayP3Model:
		im = colour.ToDisplayP3(im)
	case colour.AdobeRGBModel:
		im = colour.ToAdobeRGB(im)
	case colour.UnknownModel, colour.SRGBModel:
		// pass
	default:
		slog.Warn("Unknown or unsupported colour model", "model", model)
	}

	return im
}
