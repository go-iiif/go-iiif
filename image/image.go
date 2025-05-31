package image

import (
	"bytes"
	"context"
	"fmt"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"time"

	"github.com/aaronland/go-image/v2/colour"
	"github.com/aaronland/go-image/v2/decode"
	"github.com/aaronland/go-image/v2/encode"
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

	ctx := context.Background()

	logger.Info("WTF", "len", len(im.Body()), "type", im.ContentType())

	im_buf := bytes.NewReader(im.Body())
	goimg, _, _, err := decode.DecodeImage(ctx, im_buf)

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

	ctx := context.Background()
	wr := new(bytes.Buffer)
	var err error

	switch content_type {
	case "jpg", "jpeg", "image/jpeg":
		err = encode.EncodeJPEG(ctx, wr, goimg, nil, nil)
	case "png", "image/png":
		err = encode.EncodePNG(ctx, wr, goimg, nil)
	case "tiff", "image/tiff":
		err = encode.EncodeTIFF(ctx, wr, goimg, nil, nil)
	case "bmp", "image/bmp":
		err = encode.EncodeBMP(ctx, wr, goimg, nil)
	case "heic", "image/heic":
		err = encode.EncodeHEIC(ctx, wr, goimg, nil)
	case "gif", "image/gif":
		err = encode.EncodeGIF(ctx, wr, goimg, nil, nil)
	default:
		err = fmt.Errorf("Unsupported filetype (%s)", content_type)
	}

	if err != nil {
		return nil, err
	}

	return wr.Bytes(), nil
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
