// Package prism provides a set of tools for colour management and conversion.
// Subpackages provide support for encoding/decoding image pixel data in
// specific colour spaces, and conversions between those spaces.
package prism

import (
	"github.com/mandykoh/go-parallel"
	"image"
	"image/color"
	"image/draw"
)

// ConvertImageToNRGBA is a convenience function for getting an NRGBA image from
// any image. If the specified image isn’t already NRGBA, a conversion is
// performed.
//
// parallelism specifies the maximum degree of parallel processing; a value of 4
// indicates that processing may be spread across up to four threads. However,
// this is not guaranteed as not all conversions are parallelised.
func ConvertImageToNRGBA(img image.Image, parallelism int) *image.NRGBA {
	switch inputImg := img.(type) {

	case *image.NRGBA:
		return inputImg

	case *image.YCbCr:
		outputImg := image.NewNRGBA(inputImg.Rect)

		parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
			for i := outputImg.Rect.Min.Y + workerNum; i < outputImg.Rect.Max.Y; i += workerCount {
				for j := outputImg.Rect.Min.X; j < outputImg.Rect.Max.X; j++ {
					c := inputImg.YCbCrAt(j, i)
					r, g, b := color.YCbCrToRGB(c.Y, c.Cb, c.Cr)
					nrgba := color.NRGBA{R: r, G: g, B: b, A: 255}
					outputImg.SetNRGBA(j, i, nrgba)
				}
			}
		})

		return outputImg

	default:
		outputImg := image.NewNRGBA(img.Bounds())
		draw.Draw(outputImg, outputImg.Rect, img, outputImg.Rect.Min, draw.Src)
		return outputImg
	}
}

// ConvertImageToRGBA is a convenience function for getting an RGBA image
// from any image. If the specified image isn’t already RGBA, a conversion is
// performed.
//
// parallelism specifies the maximum degree of parallel processing; a value of 4
// indicates that processing may be spread across up to four threads. However,
// this is not guaranteed as not all conversions are parallelised.
func ConvertImageToRGBA(img image.Image, parallelism int) *image.RGBA {
	switch inputImg := img.(type) {

	case *image.RGBA:
		return inputImg

	case *image.RGBA64:
		outputImg := image.NewRGBA(inputImg.Rect)

		parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
			for i := outputImg.Rect.Min.Y + workerNum; i < outputImg.Rect.Max.Y; i += workerCount {
				for j := outputImg.Rect.Min.X; j < outputImg.Rect.Max.X; j++ {
					inputOffset := inputImg.PixOffset(j, i)
					outputOffset := outputImg.PixOffset(j, i)
					outputImg.Pix[outputOffset] = inputImg.Pix[inputOffset]
					outputImg.Pix[outputOffset+1] = inputImg.Pix[inputOffset+2]
					outputImg.Pix[outputOffset+2] = inputImg.Pix[inputOffset+4]
					outputImg.Pix[outputOffset+3] = inputImg.Pix[inputOffset+6]
				}
			}
		})

		return outputImg

	default:
		outputImg := image.NewRGBA(img.Bounds())
		draw.Draw(outputImg, outputImg.Rect, img, outputImg.Rect.Min, draw.Src)
		return outputImg
	}
}

// ConvertImageToRGBA64 is a convenience function for getting an RGBA64 image
// from any image. If the specified image isn’t already RGBA64, a conversion is
// performed.
//
// parallelism specifies the maximum degree of parallel processing; a value of 4
// indicates that processing may be spread across up to four threads. However,
// this is not guaranteed as not all conversions are parallelised.
func ConvertImageToRGBA64(img image.Image, parallelism int) *image.RGBA64 {
	switch inputImg := img.(type) {

	case *image.NRGBA:
		outputImg := image.NewRGBA64(inputImg.Rect)

		parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
			for i := outputImg.Rect.Min.Y + workerNum; i < outputImg.Rect.Max.Y; i += workerCount {
				for j := outputImg.Rect.Min.X; j < outputImg.Rect.Max.X; j++ {
					r, g, b, a := inputImg.NRGBAAt(j, i).RGBA()
					rgba64 := color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}
					outputImg.SetRGBA64(j, i, rgba64)
				}
			}
		})
		return outputImg

	case *image.RGBA:
		outputImg := image.NewRGBA64(inputImg.Rect)

		parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
			for i := outputImg.Rect.Min.Y + workerNum; i < outputImg.Rect.Max.Y; i += workerCount {
				for j := outputImg.Rect.Min.X; j < outputImg.Rect.Max.X; j++ {
					inputOffset := inputImg.PixOffset(j, i)
					outputOffset := outputImg.PixOffset(j, i)
					outputImg.Pix[outputOffset] = inputImg.Pix[inputOffset]
					outputImg.Pix[outputOffset+1] = inputImg.Pix[inputOffset]
					outputImg.Pix[outputOffset+2] = inputImg.Pix[inputOffset+1]
					outputImg.Pix[outputOffset+3] = inputImg.Pix[inputOffset+1]
					outputImg.Pix[outputOffset+4] = inputImg.Pix[inputOffset+2]
					outputImg.Pix[outputOffset+5] = inputImg.Pix[inputOffset+2]
					outputImg.Pix[outputOffset+6] = inputImg.Pix[inputOffset+3]
					outputImg.Pix[outputOffset+7] = inputImg.Pix[inputOffset+3]
				}
			}
		})

		return outputImg

	case *image.RGBA64:
		return inputImg

	case *image.YCbCr:
		outputImg := image.NewRGBA64(inputImg.Rect)

		parallel.RunWorkers(parallelism, func(workerNum, workerCount int) {
			for i := outputImg.Rect.Min.Y + workerNum; i < outputImg.Rect.Max.Y; i += workerCount {
				for j := outputImg.Rect.Min.X; j < outputImg.Rect.Max.X; j++ {
					r, g, b, _ := inputImg.YCbCrAt(j, i).RGBA()
					rgba64 := color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: 65535}
					outputImg.SetRGBA64(j, i, rgba64)
				}
			}
		})

		return outputImg

	default:
		outputImg := image.NewRGBA64(img.Bounds())
		draw.Draw(outputImg, outputImg.Rect, img, outputImg.Rect.Min, draw.Src)
		return outputImg
	}
}
