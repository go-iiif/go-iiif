// Package rotate provides methods for rotating images.
package rotate

import (
	"context"
	"fmt"
	"image"

	"github.com/aaronland/go-image/imaging"
)

// RotateImageWithOrientation will rotate 'im' based on EXIF orientation value defined in 'orientation'.
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

// RotateImageWithOrientation will rotate 'im' by 'degrees' degrees. Currently only values in units of
// 90 degrees or 0 are supported.
func RotateImageWithDegrees(ctx context.Context, im image.Image, degrees float64) (image.Image, error) {

	// See also: https://github.com/anthonynsimon/bild#rotate
	// The problem is that bild doesn't rotate the "canvas" just the image

	switch degrees {
	case 90.0:
		im = imaging.Rotate90(im)
	case 180.0:
		im = imaging.Rotate180(im)
	case 270.0:
		im = imaging.Rotate270(im)
	default:
		return nil, fmt.Errorf("Unsupported value, %f", degrees)
	}

	return im, nil
}
