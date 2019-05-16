package noborders

import (
	"errors"
	"image"
	"image/draw"
)

// RemoveBorders removes low-entropy/low-variance border around the specified image.
func RemoveBorders(img image.Image, opts Options) (result image.Image, err error) {
	if img == nil {
		return result, errors.New("image must not be nil")
	}

	if opts == nil {
		opts = Opts()
	}

	var bounds = img.Bounds()
	var rounds int
	var quiescent bool
	for !quiescent {
		rounds++
		var lastBounds = bounds
		rowEntropies, colEntropies := sliceOperation(img, bounds, entropy)
		rowVariances, colVariances := sliceOperation(img, bounds, variance)
		// rows
		for i := 0; i < len(rowEntropies); i++ {
			if rowEntropies[i] < opts.Entropy() || rowVariances[i] < opts.Variance() {
				bounds.Min.Y++
			} else {
				break
			}
		}
		for i := len(rowEntropies) - 1; i >= 0; i-- {
			if rowEntropies[i] < opts.Entropy() || rowVariances[i] < opts.Variance() {
				bounds.Max.Y--
			} else {
				break
			}
		}
		if bounds.Max.Y < bounds.Min.Y {
			bounds.Max.Y = bounds.Min.Y
		}

		// cols
		for i := 0; i < len(colEntropies); i++ {
			if colEntropies[i] < opts.Entropy() || colVariances[i] < opts.Variance() {
				bounds.Min.X++
			} else {
				break
			}
		}
		for i := len(colEntropies) - 1; i >= 0; i-- {
			if colEntropies[i] < opts.Entropy() || colVariances[i] < opts.Variance() {
				bounds.Max.X--
			} else {
				break
			}
		}
		if bounds.Max.X < bounds.Min.X {
			bounds.Max.X = bounds.Min.X
		}

		if bounds == lastBounds || !opts.MultiPass() {
			quiescent = true
		}
	}

	result = cropImage(img, bounds)
	return
}

// cropImage crops an image.
func cropImage(img image.Image, crop image.Rectangle) image.Image {
	result := image.NewRGBA(crop)
	draw.Draw(result, crop, img, crop.Min, draw.Src)
	return result
}
