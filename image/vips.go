package image

// https://github.com/h2non/bimg
// https://github.com/jcupitt/libvips

import (
	"bytes"
	"errors"
	"fmt"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
	"gopkg.in/h2non/bimg.v1"
	"image"
	"image/gif"
	_ "log"
	"strconv"
	"strings"
)

type VIPSImage struct {
	Image
	config    *iiifconfig.Config
	source    iiifsource.Source
	source_id string
	id        string
	bimg      *bimg.Image
	isgif     bool
}

type VIPSDimensions struct {
	Dimensions
	imagesize bimg.ImageSize
}

func (d *VIPSDimensions) Height() int {
	return d.imagesize.Height
}

func (d *VIPSDimensions) Width() int {
	return d.imagesize.Width
}

/*

See notes in NewVIPSImageFromConfigWithSource - basically getting an image's
dimensions after the we've done the GIF conversion (just see the notes...)
will make bimg/libvips sad so account for that in Dimensions() and create a
pure Go implementation of the Dimensions interface (20160922/thisisaaronland)

*/

type GolangImageDimensions struct {
	Dimensions
	image image.Image
}

func (dims *GolangImageDimensions) Width() int {
	bounds := dims.image.Bounds()
	return bounds.Max.X
}

func (dims *GolangImageDimensions) Height() int {
	bounds := dims.image.Bounds()
	return bounds.Max.Y
}

func NewVIPSImageFromConfigWithSource(config *iiifconfig.Config, src iiifsource.Source, id string) (*VIPSImage, error) {

	body, err := src.Read(id)

	if err != nil {
		return nil, err
	}

	bimg := bimg.NewImage(body)

	im := VIPSImage{
		config:    config,
		source:    src,
		source_id: id,
		id:        id,
		bimg:      bimg,
		isgif:     false,
	}

	/*

		Hey look - see the 'isgif' flag? We're going to hijack the fact that
		bimg doesn't handle GIF files and if someone requests them then we
		will do the conversion after the final call to im.bimg.Process and
		after we do handle any custom features. We are relying on the fact
		that both bimg.NewImage and bimg.Image() expect and return raw bytes
		and we are ignoring whatever bimg thinks in the Format() function.
		So basically you should not try to any processing in bimg/libvips
		after the -> GIF transformation. (20160922/thisisaaronland)

		See also: https://github.com/h2non/bimg/issues/41
	*/

	return &im, nil
}

func (im *VIPSImage) Update(body []byte) error {

	bimg := bimg.NewImage(body)
	im.bimg = bimg

	return nil
}

func (im *VIPSImage) Body() []byte {

	return im.bimg.Image()
}

func (im *VIPSImage) Format() string {

	// see notes in NewVIPSImageFromConfigWithSource

	if im.isgif {
		return "gif"
	}

	return im.bimg.Type()
}

func (im *VIPSImage) ContentType() string {

	format := im.Format()

	if format == "jpg" || format == "jpeg" {
		return "image/jpeg"
	} else if format == "png" {
		return "image/png"
	} else if format == "webp" {
		return "image/webp"
	} else if format == "tif" || format == "tiff" {
		return "image/tiff"
	} else if format == "gif" {
		return "image/gif"
	} else {
		return ""
	}
}

func (im *VIPSImage) Identifier() string {
	return im.id
}

func (im *VIPSImage) Rename(id string) error {
	im.id = id
	return nil
}

func (im *VIPSImage) Dimensions() (Dimensions, error) {

	// see notes in NewVIPSImageFromConfigWithSource
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

		/*
		   So here's a thing that we need to do because... computers?
		   (20160910/thisisaaronland)
		*/

		if opts.Top == 0 && opts.Left == 0 {
			opts.Top = -1
		}

		_, err = im.bimg.Process(opts)

		if err != nil {
			return err
		}

	}

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

	fi, err := t.FormatInstructions(im)

	if err != nil {
		return nil
	}

	if fi.Format == "jpg" {
		opts.Type = bimg.JPEG
	} else if fi.Format == "png" {
		opts.Type = bimg.PNG
	} else if fi.Format == "webp" {
		opts.Type = bimg.WEBP
	} else if fi.Format == "tif" {
		opts.Type = bimg.TIFF
	} else if fi.Format == "gif" {
		opts.Type = bimg.PNG // see this - we're just going to trick libvips until the very last minute...
	} else {
		msg := fmt.Sprintf("Unsupported image format '%s'", fi.Format)
		return errors.New(msg)
	}

	_, err = im.bimg.Process(opts)

	if err != nil {
		return err
	}

	// None of what follows is part of the IIIF spec so it's not clear
	// to me yet how to make this in to a sane interface. For the time
	// being since there is only lipvips we'll just take the opportunity
	// to think about it... (20160917/thisisaaronland)

	// Also note the way we are diligently setting in `im.isgif` in each
	// of the features below. That's because this is a bimg/libvips-ism
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

	// see notes in NewVIPSImageFromConfigWithSource

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
