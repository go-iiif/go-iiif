package noborders

import (
	"image"
	"image/color"
	"math"

	"github.com/gonum/stat"
)

// sliceOperation performs a transform operation on every row and column slice on the
// specified image and crop.
func sliceOperation(img image.Image, crop image.Rectangle, oper func(img image.Image, r image.Rectangle) float64) (rowResults, colResults []float64) {
	for x := crop.Min.X; x < crop.Max.X; x++ {
		var col = image.Rect(x, crop.Min.Y, x+1, crop.Max.Y)
		colResults = append(colResults, oper(img, col))
	}
	for y := crop.Min.Y; y < crop.Max.Y; y++ {
		var row = image.Rect(crop.Min.X, y, crop.Max.X, y+1)
		rowResults = append(rowResults, oper(img, row))
	}

	return
}

// variance computes the pixel intensity variance of a portion of an image.
func variance(img image.Image, r image.Rectangle) float64 {
	var vals []float64
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			vals = append(vals, float64(greyvalue(img.At(x, y))))
		}
	}

	return stat.Variance(vals, nil)
}

// entropy calculates the entropy of a portion of an image.
// Adapted from https://github.com/iand/salience (modified to handle 1-pixel image slices)
// Who adapted it from http://www.astro.cornell.edu/research/projects/compression/entropy.html
func entropy(img image.Image, r image.Rectangle) float64 {
	arraySize := 256*2 - 1
	freq := make([]float64, arraySize)

	if r.Max.X-r.Min.X < 2 && r.Max.Y-r.Min.Y < 2 {
		return 0.0
	}

	if r.Max.Y-r.Min.Y < 2 {
		for x := r.Min.X; x < r.Max.X-1; x++ {
			for y := r.Min.Y; y < r.Max.Y; y++ {
				diff := greyvalue(img.At(x, y)) - greyvalue(img.At(x+1, y))
				if -(arraySize+1)/2 < diff && diff < (arraySize+1)/2 {
					freq[diff+(arraySize-1)/2]++
				}
			}
		}
	} else {
		for y := r.Min.Y; y < r.Max.Y-1; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				diff := greyvalue(img.At(x, y)) - greyvalue(img.At(x, y+1))
				if -(arraySize+1)/2 < diff && diff < (arraySize+1)/2 {
					freq[diff+(arraySize-1)/2]++
				}
			}
		}
	}

	n := 0.0
	for _, v := range freq {
		n += v
	}

	e := 0.0
	for i := 0; i < len(freq); i++ {
		freq[i] = freq[i] / n
		if freq[i] != 0.0 {
			e -= freq[i] * math.Log2(freq[i])
		}
	}

	return e

}

// greyvalue computes the greyscale value of a Color based on the luminosity method.
func greyvalue(c color.Color) int {
	r, g, b, _ := c.RGBA()
	return int((r*299 + g*587 + b*114) / 1000)
}
