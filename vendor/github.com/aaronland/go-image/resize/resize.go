// Package resize provides methods for resizing images.
package resize

import (
	"context"
	"image"
	"math"

	nfnt_resize "github.com/nfnt/resize"
)

// ResizeImages scaled and resizes 'im' return a new `image.Image` instance whose maximum dimension is 'max'.
func ResizeImage(ctx context.Context, im image.Image, max int) (image.Image, error) {

	// calculating w,h is probably unnecessary since we're
	// calling resize.Thumbnail but it will do for now...
	// (20180708/thisisaaronland)

	bounds := im.Bounds()
	dims := bounds.Max

	width := dims.X
	height := dims.Y

	ratio_w := float64(max) / float64(width)
	ratio_h := float64(max) / float64(height)

	ratio := math.Min(ratio_w, ratio_h)

	w := uint(float64(width) * ratio)
	h := uint(float64(height) * ratio)

	sm := nfnt_resize.Thumbnail(w, h, im, nfnt_resize.Lanczos3)

	return sm, nil
}
