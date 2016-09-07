package image

// https://github.com/h2non/bimg
// https://github.com/jcupitt/libvips

import (
	"fmt"
	"github.com/thisisaaronland/go-iiif/source"
	"gopkg.in/h2non/bimg.v1"
)

type VIPSImage struct {
	Image
	source source.Source
	id     string
	bimg   *bimg.Image
}

type VIPSDimensions struct {
	Dimensions
	imagesize bimg.ImageSize
}

func NewVIPSImageFromSource(src source.Source, id string) (*VIPSImage, error) {

	body, err := src.Read(id)

	if err != nil {
		return nil, err
	}

	bimg := bimg.NewImage(body)

	im := VIPSImage{
		source: src,
		id:     id,
		bimg:   bimg,
	}

	return &im, nil
}

func (im *VIPSImage) Body() []byte {
	return im.bimg.Image()
}

func (im *VIPSImage) Format() string {
	return im.bimg.Type()
}

func (im *VIPSImage) ContentType() string {

	format := im.Format()

	if format == "jpg" || format == "jpeg" {
		return "image/jpg"
	} else if format == "png" {
		return "image/png"
	} else if format == "webp" {
		return "image/webp"
	} else if format == "tif" || format == "tiff" {
		return "image/tiff"
	} else {
		return ""
	}
}

func (im *VIPSImage) Identifier() string {
	return im.id
}

func (im *VIPSImage) Dimensions() (Dimensions, error) {

	sz, err := im.bimg.Size()

	if err != nil {
		return nil, err
	}

	d := VIPSDimensions{
		imagesize: sz,
	}

	return &d, nil
}

// https://godoc.org/github.com/h2non/bimg#Options

func (im *VIPSImage) Transform(t *Transformation) error {

	var opts bimg.Options

	if t.Region != "full" {

		rgi, err := t.RegionInstructions(im)

		if err != nil {
			return err
		}

		opts = bimg.Options{
			AreaWidth:  rgi.Width,
			AreaHeight: rgi.Height,
			Left:       rgi.X,
			Top:        rgi.Y,
		}

		_, err = im.bimg.Process(opts)

		if err != nil {
			return err
		}
	}

	// QUESTION: after Process what are the dimensions of im.bimg ?

	dims, err := im.Dimensions()

	if err != nil {
		return err
	}

	opts = bimg.Options{
		Width:  dims.Width(),  // opts.AreaWidth,
		Height: dims.Height(), // opts.AreaHeight,
	}

	if t.Size != "max" && t.Size != "full" {

		si, err := t.SizeInstructions(im)

		if err != nil {
			return err
		}

		opts.Height = si.Height
		opts.Width = si.Width
		opts.Enlarge = si.Enlarge
		opts.Force = si.Force
	}

	ri, err := t.RotationInstructions(im)

	if err != nil {
		return nil
	}

	// PLEASE FIX ME
	// var rotationMissing = "libvips cannot rotate angle that isn't a multiple of 90: %#v"

	opts.Flip = ri.Flip
	opts.Rotate = bimg.Angle(ri.Angle % 360)

	if t.Quality == "color" || t.Quality == "default" {
		// do nothing.
	} else if t.Quality == "gray" {
		opts.Interpretation = bimg.InterpretationBW
	} else if t.Quality == "bitonal" {
		opts.Interpretation = bimg.InterpretationBW
	} else {
		// this should be trapped above
	}

	_, err = im.bimg.Process(opts)

	if err != nil {
		return err
	}

	return nil
}

func (d *VIPSDimensions) Height() int {
	return d.imagesize.Height
}

func (d *VIPSDimensions) Width() int {
	return d.imagesize.Width
}
