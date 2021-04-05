package image

import (
	"github.com/anthonynsimon/bild/effect"
)

// "spork": { "syntax": "spork", "required": false, "supported": true, "match": "^sharpen$" }

func SharpenImage(im Image) error {

	goimg, err := IIIFImageToGolangImage(im)

	if err != nil {
		return err
	}

	sharpened := effect.Sharpen(goimg)

	return GolangImageToIIIFImage(sharpened, im)
}
