package transformation

import (
	iiifimage "github.com/go-iiif/go-iiif/v5/image"	
	"github.com/anthonynsimon/bild/effect"
	_ "log"
)

type CrispImageOptions struct {
	Radius float64
	Amount float64
	Median float64
}

func DefaultCrispImageOptions() *CrispImageOptions {

	opts := &CrispImageOptions{
		Radius: 2.0,
		Amount: 0.5,
		Median: 0.025,
	}

	return opts
}

func CrispImage(im iiifimage.Image, opts *CrispImageOptions) error {

	new_im, err := iiifimage.IIIFImageToGolangImage(im)

	if err != nil {
		return err
	}

	new_im = effect.UnsharpMask(new_im, opts.Radius, opts.Amount)
	new_im = effect.Median(new_im, opts.Median)

	return iiifimage.GolangImageToIIIFImage(new_im, im)
}
