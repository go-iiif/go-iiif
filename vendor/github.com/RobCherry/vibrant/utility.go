package vibrant

import (
	"golang.org/x/image/draw"
	"image"
	"math"
)

// Utility function for restricting the value of a number.
func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func clampUint8(value, min, max uint8) uint8 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Utility function for rounding the value of a number.
func roundFloat64(value float64) float64 {
	if value < 0.0 {
		return value - 0.5
	}
	return value + 0.5
}

// SubImager is a utility interface for an image.Image that can extract a sub-image.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// ScaleImageDown will scale the image down as needed.  This is relatively slow.
func ScaleImageDown(src image.Image, resizeArea uint64, scaler draw.Scaler) image.Image {
	if resizeArea > 0 && scaler != nil {
		imageSize := src.Bounds().Size()
		imageArea := uint64(imageSize.X * imageSize.Y)
		if imageArea > resizeArea {
			scaleRatio := math.Sqrt(float64(resizeArea) / float64(imageArea))
			scaled := image.NewNRGBA(image.Rect(0, 0, int(math.Floor(float64(imageSize.X)*scaleRatio)), int(math.Floor(float64(imageSize.Y)*scaleRatio))))
			scaler.Scale(scaled, scaled.Bounds(), src, src.Bounds(), draw.Src, nil)
			return scaled
		}
	}
	// Scaling has been disabled or is not needed
	return src
}
