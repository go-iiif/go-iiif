package native

import (
	"bytes"
	"errors"
	_ "fmt"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiifsource "github.com/go-iiif/go-iiif/source"
	"image"
	"image/gif"
	_ "log"
	"strconv"
	"strings"
)

type NativeImage struct {
	iiifimage.Image
	config    *iiifconfig.Config
	source    iiifsource.Source
	source_id string
	id        string
	img       image.Image
	isgif     bool
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

	img := img.NewImage(body)
	im.img = img

	return nil
}

func (im *NativeImage) Body() []byte {

	return im.img.Image()
}

func (im *NativeImage) Format() string {

	return im.img.Type()
}

func (im *NativeImage) ContentType() string {

	format := im.Format()

	if format == "jpg" || format == "jpeg" {
		return "image/jpeg"
	} else if format == "png" {
		return "image/png"
	} else if format == "webp" {
		return "image/webp"
	} else if format == "svg" {
		return "image/svg+xml"
	} else if format == "tif" || format == "tiff" {
		return "image/tiff"
	} else if format == "gif" {
		return "image/gif"
	} else {
		return ""
	}
}

func (im *NativeImage) Identifier() string {
	return im.id
}

func (im *NativeImage) Rename(id string) error {
	im.id = id
	return nil
}

func (im *NativeImage) Dimensions() (Dimensions, error) {

	// see notes in NewNativeImageFromConfigWithSource
	// ideally this never gets triggered but just in case...

	if im.isgif {

		buf := bytes.NewBuffer(im.Body())
		goimg, err := gif.Decode(buf)

		if err != nil {
			return nil, err
		}

		d := GolangImageDimensions{
			image: goimg,
		}

		return &d, nil
	}

	sz, err := im.img.Size()

	if err != nil {
		return nil, err
	}

	d := NativeDimensions{
		imagesize: sz,
	}

	return &d, nil
}

// https://godoc.org/github.com/h2non/img#Options

func (im *NativeImage) Transform(t *Transformation) error {

	// ... PLEASE WRITE ME ....

	// PLEASE PUT THIS IN A COMMON PACKAGE

	// None of what follows is part of the IIIF spec so it's not clear
	// to me yet how to make this in to a sane interface. For the time
	// being since there is only lipvips we'll just take the opportunity
	// to think about it... (20160917/thisisaaronland)

	// Also note the way we are diligently setting in `im.isgif` in each
	// of the features below. That's because this is a img/libvips-ism
	// and we assume that any of these can encode GIFs because pure-Go and
	// the rest of the code does need to know about it...
	// (20160922/thisisaaronland)

	if t.Quality == "dither" {

		err = DitherImage(im)

		if err != nil {
			return err
		}

		if fi.Format == "gif" {
			im.isgif = true
		}

	} else if strings.HasPrefix(t.Quality, "primitive:") {

		parts := strings.Split(t.Quality, ":")
		parts = strings.Split(parts[1], ",")

		mode, err := strconv.Atoi(parts[0])

		if err != nil {
			return err
		}

		iters, err := strconv.Atoi(parts[1])

		if err != nil {
			return err
		}

		max_iters := im.config.Primitive.MaxIterations

		if max_iters > 0 && iters > max_iters {
			return errors.New("Invalid primitive iterations")
		}

		alpha, err := strconv.Atoi(parts[2])

		if err != nil {
			return err
		}

		if alpha > 255 {
			return errors.New("Invalid primitive alpha")
		}

		animated := false

		if fi.Format == "gif" {
			animated = true
		}

		opts := PrimitiveOptions{
			Alpha:      alpha,
			Mode:       mode,
			Iterations: iters,
			Size:       0,
			Animated:   animated,
		}

		err = PrimitiveImage(im, opts)

		if err != nil {
			return err
		}

		if fi.Format == "gif" {
			im.isgif = true
		}
	}

	// END OF none of what follows is part of the IIIF spec

	// see notes in NewNativeImageFromConfigWithSource

	if fi.Format == "gif" && !im.isgif {

		goimg, err := IIIFImageToGolangImage(im)

		if err != nil {
			return err
		}

		im.isgif = true

		err = GolangImageToIIIFImage(goimg, im)

		if err != nil {
			return err
		}

	}

	return nil
}
