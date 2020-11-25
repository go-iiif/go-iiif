package rotate

import (
	"context"
	"errors"
	"github.com/aaronland/go-image-rotate/imaging"
	"image"
)

func RotateImageWithOrientation(ctx context.Context, im image.Image, orientation string) (image.Image, error) {

	switch orientation {
	case "1":
		// pass
	case "2":
		im = imaging.FlipV(im)
	case "3":
		im = imaging.Rotate180(im)
	case "4":
		im = imaging.Rotate180(imaging.FlipV(im))
	case "5":
		im = imaging.Rotate270(imaging.FlipV(im))
	case "6":
		im = imaging.Rotate270(im)
	case "7":
		im = imaging.Rotate90(imaging.FlipV(im))
	case "8":
		im = imaging.Rotate90(im)
	}

	return im, nil
}

func RotateImageWithDegrees(ctx context.Context, im image.Image, degrees float64) (image.Image, error) {

	switch degrees {
	case 90.0:
		im = imaging.Rotate90(im)
	case 180.0:
		im = imaging.Rotate180(im)
	case 270.0:
		im = imaging.Rotate270(im)
	default:
		return nil, errors.New("Unsupported value")
	}

	return im, nil
}
