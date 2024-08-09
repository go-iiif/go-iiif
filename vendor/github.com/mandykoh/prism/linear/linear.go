package linear

import (
	"github.com/mandykoh/go-parallel"
	"image"
	"image/color"
	"image/draw"
)

// NormalisedTo8Bit clamps and scales a normalised value to the range 0-255.
func NormalisedTo8Bit(v float32) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 255
	}
	return uint8(v*255 + 0.5)
}

// NormalisedTo9Bit clamps and scales a normalised value to the range 0-511.
func NormalisedTo9Bit(v float32) uint16 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 511
	}
	return uint16(v*511 + 0.5)
}

// NormalisedTo16Bit clamps and scales a normalised value to the range 0-65535.
func NormalisedTo16Bit(v float32) uint16 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 65535
	}
	return uint16(v*65535 + 0.5)
}

// TransformImageColor applies a colour transformation function to all pixels of
// src, writing the results to dst at its origin.
//
// src and dst may be the same image.
//
// parallelism specifies the degree of parallel processing; a value of 4
// indicates that processing will be spread across four threads.
func TransformImageColor(dst draw.Image, src image.Image, parallelism int, transformColor func(color.Color) color.RGBA64) {
	bounds := src.Bounds()
	dstOffsetX := dst.Bounds().Min.X - bounds.Min.X
	dstOffsetY := dst.Bounds().Min.Y - bounds.Min.Y

	switch dstImg := dst.(type) {

	case *image.RGBA64:
		if srcImg, ok := src.(*image.RGBA64); ok {

			parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
				for i := bounds.Min.Y + workerNum; i < bounds.Max.Y; i += workerCount {
					for j := bounds.Min.X; j < bounds.Max.X; j++ {
						c := transformColor(srcImg.RGBA64At(j, i))

						offset := dstImg.PixOffset(j+dstOffsetX, i+dstOffsetY)
						dstImg.Pix[offset] = uint8(c.R >> 8)
						dstImg.Pix[offset+1] = uint8(c.R & 0xFF)
						dstImg.Pix[offset+2] = uint8(c.G >> 8)
						dstImg.Pix[offset+3] = uint8(c.G & 0xFF)
						dstImg.Pix[offset+4] = uint8(c.B >> 8)
						dstImg.Pix[offset+5] = uint8(c.B & 0xFF)
						dstImg.Pix[offset+6] = uint8(c.A >> 8)
						dstImg.Pix[offset+7] = uint8(c.A & 0xFF)
					}
				}
			})

		} else {
			parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
				for i := bounds.Min.Y + workerNum; i < bounds.Max.Y; i += workerCount {
					for j := bounds.Min.X; j < bounds.Max.X; j++ {
						c := transformColor(src.At(j, i))

						offset := dstImg.PixOffset(j+dstOffsetX, i+dstOffsetY)
						dstImg.Pix[offset] = uint8(c.R >> 8)
						dstImg.Pix[offset+1] = uint8(c.R & 0xFF)
						dstImg.Pix[offset+2] = uint8(c.G >> 8)
						dstImg.Pix[offset+3] = uint8(c.G & 0xFF)
						dstImg.Pix[offset+4] = uint8(c.B >> 8)
						dstImg.Pix[offset+5] = uint8(c.B & 0xFF)
						dstImg.Pix[offset+6] = uint8(c.A >> 8)
						dstImg.Pix[offset+7] = uint8(c.A & 0xFF)
					}
				}
			})
		}

	case *image.RGBA:
		parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
			for i := bounds.Min.Y + workerNum; i < bounds.Max.Y; i += workerCount {
				for j := bounds.Min.X; j < bounds.Max.X; j++ {
					c := transformColor(src.At(j, i))

					offset := dstImg.PixOffset(j+dstOffsetX, i+dstOffsetY)
					dstImg.Pix[offset] = uint8(c.R >> 8)
					dstImg.Pix[offset+1] = uint8(c.G >> 8)
					dstImg.Pix[offset+2] = uint8(c.B >> 8)
					dstImg.Pix[offset+3] = uint8(c.A >> 8)
				}
			}
		})

	default:
		parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
			for i := bounds.Min.Y + workerNum; i < bounds.Max.Y; i += workerCount {
				for j := bounds.Min.X; j < bounds.Max.X; j++ {
					dst.Set(j+dstOffsetX, i+dstOffsetY, transformColor(src.At(j, i)))
				}
			}
		})
	}
}
