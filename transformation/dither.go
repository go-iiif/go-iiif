package transformation

import (
	iiifimage "github.com/go-iiif/go-iiif/v5/image"	
	"github.com/koyachi/go-atkinson"
)

func DitherImage(im iiifimage.Image) error {

	goimg, err := iiifimage.IIIFImageToGolangImage(im)

	if err != nil {
		return err
	}

	dithered, err := atkinson.Dither(goimg)

	if err != nil {
		return err
	}

	return iiifimage.GolangImageToIIIFImage(dithered, im)
}
