package extruder

import (
	"errors"
	"github.com/aaronland/go-colours"
	"github.com/neocortical/noborders"
	"github.com/nfnt/resize"
	"image"
	"strings"
)

func NewNamedExtruder(name string, args ...interface{}) (colours.Extruder, error) {

	var ex colours.Extruder
	var err error

	switch strings.ToUpper(name) {
	case "SIMPLE":
		ex, err = NewSimpleExtruder(args...)
	case "VIBRANT":
		ex, err = NewVibrantExtruder(args...)
	default:
		err = errors.New("Invalid or unknown extruder")
	}

	return ex, err
}

func PrepareImage(im image.Image) (image.Image, error) {

	im = resize.Resize(100, 0, im, resize.Bilinear)

	opts := noborders.Opts()
	opts.SetEntropy(0.05)
	opts.SetVariance(100000)
	opts.SetMultiPass(true)

	return noborders.RemoveBorders(im, opts)
}
