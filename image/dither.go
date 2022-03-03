package image

import (
	"github.com/koyachi/go-atkinson"
)

func DitherImage(im Image) error {

	goimg, err := IIIFImageToGolangImage(im)

	if err != nil {
		return err
	}

	dithered, err := atkinson.Dither(goimg)

	if err != nil {
		return err
	}

	return GolangImageToIIIFImage(dithered, im)
}
