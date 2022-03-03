package image

import (
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

func CrispImage(im Image, opts *CrispImageOptions) error {

	new_im, err := IIIFImageToGolangImage(im)

	if err != nil {
		return err
	}

	new_im = effect.UnsharpMask(new_im, opts.Radius, opts.Amount)
	new_im = effect.Median(new_im, opts.Median)

	return GolangImageToIIIFImage(new_im, im)
}
