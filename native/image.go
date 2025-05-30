package native

// to consider... https://github.com/corona10/goimghdr

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"log/slog"

	"github.com/aaronland/go-image/v2/colour"
	"github.com/aaronland/go-image/v2/decode"
	"github.com/aaronland/go-image/v2/encode"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/segment"
	"github.com/anthonynsimon/bild/transform"
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiifimage "github.com/go-iiif/go-iiif/v8/image"
	iiifsource "github.com/go-iiif/go-iiif/v8/source"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
)

type NativeImage struct {
	iiifimage.Image
	config    *iiifconfig.Config
	source    iiifsource.Source
	source_id string
	id        string
	img       image.Image
	format    string
	model     colour.Model
}

type NativeDimensions struct {
	iiifimage.Dimensions
	bounds image.Rectangle
}

func (d *NativeDimensions) Height() int {
	return d.bounds.Max.Y
}

func (d *NativeDimensions) Width() int {
	return d.bounds.Max.X
}

func (im *NativeImage) Update(body []byte) error {

	img, img_fmt, err := decodeImageBytes(body)

	if err != nil {
		return fmt.Errorf("Failed to decode image bytes, %w", err)
	}

	im.img = img
	im.format = img_fmt

	return nil
}

func (im *NativeImage) Body() []byte {
	body, _ := encodeImage(im.img, im.format)
	return body
}

func (im *NativeImage) Format() string {

	return im.format
}

func (im *NativeImage) ColourModel() colour.Model {

	return im.model
}

func (im *NativeImage) ContentType() string {
	return im.Format()
}

func (im *NativeImage) Identifier() string {
	return im.id
}

func (im *NativeImage) Rename(id string) error {
	im.id = id
	return nil
}

func (im *NativeImage) Dimensions() (iiifimage.Dimensions, error) {

	dims := &NativeDimensions{
		bounds: im.img.Bounds(),
	}

	return dims, nil
}

func (im *NativeImage) Transform(t *iiifimage.Transformation) error {

	if t.Region != "full" {

		rgi, err := t.RegionInstructions(im)

		if err != nil {
			return fmt.Errorf("Failed to derive region instructions, %w", err)
		}

		si, err := t.SizeInstructions(im)

		if err != nil {
			return fmt.Errorf("Failed to derive size instructions, %w", err)
		}

		if rgi.SmartCrop {

			resizer := nfnt.NewDefaultResizer()
			analyzer := smartcrop.NewAnalyzer(resizer)

			width := si.Width
			height := si.Height

			topCrop, err := analyzer.FindBestCrop(im.img, width, height)

			if err != nil {
				return fmt.Errorf("Failed to derive best crop, %w", err)
			}

			type SubImager interface {
				SubImage(r image.Rectangle) image.Image
			}

			cropped := im.img.(SubImager).SubImage(topCrop)

			if cropped.Bounds().Dx() != width || cropped.Bounds().Dy() != height {
				cropped = resizer.Resize(cropped, uint(width), uint(height))
			}

			im.img = cropped

		} else {

			x1 := rgi.X
			y1 := rgi.Y
			x2 := x1 + rgi.Width
			y2 := y1 + rgi.Height

			bounds := image.Rect(x1, y1, x2, y2)

			img := transform.Crop(im.img, bounds)

			// Crop returns a new image which contains the intersection between the rect and
			// the image provided as params. Only the intersection is returned. If a rect larger
			// than the image is provided, no fill is done to the 'empty' area.
			// https://godoc.org/github.com/anthonynsimon/bild/transform#Crop

			im.img = img
		}
	}

	if t.Size != "max" && t.Size != "full" {

		si, err := t.SizeInstructions(im)

		if err != nil {
			return err
		}

		w := si.Width
		h := si.Height

		// https://github.com/anthonynsimon/bild/#resize-resampling-filters
		// https://godoc.org/github.com/anthonynsimon/bild/transform#ResampleFilter

		img := transform.Resize(im.img, w, h, transform.Lanczos)
		im.img = img
	}

	ri, err := t.RotationInstructions(im)

	if err != nil {
		return nil
	}

	// auto-rotate checks... are they necessary in a plain-vanilla Go context?

	if ri.Angle > 0.0 {
		angle := float64(ri.Angle)
		img := transform.Rotate(im.img, angle, nil)
		im.img = img
	}

	if ri.Flip {
		img := transform.FlipH(im.img)
		im.img = img
	}

	switch t.Quality {
	case "color", "default":
		// do nothing.
	case "gray":
		img := effect.Grayscale(im.img)
		im.img = img
	case "bitonal":
		img := segment.Threshold(im.img, 128)
		im.img = img
	default:
		// this should be trapped above
	}

	fi, err := t.FormatInstructions(im)

	if err != nil {
		return err
	}

	encode := false

	if fi.Format != im.format {

		encode = true

		// sigh... computers, amirite?

		if fi.Format == "jpg" && im.format == "image/jpeg" {
			encode = false
		}
	}

	if encode {

		body, err := encodeImage(im.img, fi.Format)

		if err != nil {
			return fmt.Errorf("Failed to encode image, %w", err)
		}

		img, img_fmt, err := decodeImageBytes(body)

		if err != nil {
			return fmt.Errorf("Failed to decode image, %w", err)
		}

		im.img = img
		im.format = img_fmt
	}

	err = iiifimage.ApplyCustomTransformations(t, im)

	if err != nil {
		return err
	}

	return nil
}

func decodeImageBytes(body []byte) (image.Image, string, error) {

	ctx := context.Background()
	br := bytes.NewReader(body)

	// There should not be any need to rotate the image since it
	// will have already been rotated in NewImageFromConfigWithSource
	decode_opts := &decode.DecodeImageOptions{
		Rotate: false,
	}

	im, content_type, _, err := decode.DecodeImageWithOptions(ctx, br, decode_opts)

	if err != nil {
		return nil, "", err
	}

	return im, content_type, nil
}

func encodeImage(im image.Image, format string) ([]byte, error) {

	ctx := context.Background()

	var b bytes.Buffer
	wr := bufio.NewWriter(&b)

	var err error

	switch format {
	case "bmp", "image/bmp":
		err = encode.EncodeBMP(ctx, wr, im, nil)
	case "jpg", "jpeg", "image/jpeg":
		err = encode.EncodeJPEG(ctx, wr, im, nil, nil)
	case "png", "image/png":
		err = encode.EncodePNG(ctx, wr, im, nil)
	case "gif", "image/gif":
		opts := gif.Options{}
		err = gif.Encode(wr, im, &opts)
	case "tiff", "image/tiff":
		err = encode.EncodeJPEG(ctx, wr, im, nil, nil)
	case "heic", "image/heic":
		err = encode.EncodeHEIC(ctx, wr, im, nil)
	default:
		err = fmt.Errorf("Unsupported encoding, '%s'", format)
	}

	if err != nil {
		slog.Error("Failed to encode image", "error", err)
		return nil, err
	}

	wr.Flush()

	return b.Bytes(), nil
}
