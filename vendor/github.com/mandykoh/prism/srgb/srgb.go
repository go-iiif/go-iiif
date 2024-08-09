package srgb

import (
	"github.com/mandykoh/prism/linear"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/mandykoh/prism/ciexyy"
)

var PrimaryRed = ciexyy.Color{X: 0.64, Y: 0.33, YY: 1}
var PrimaryGreen = ciexyy.Color{X: 0.3, Y: 0.6, YY: 1}
var PrimaryBlue = ciexyy.Color{X: 0.15, Y: 0.06, YY: 1}
var StandardWhitePoint = ciexyy.D65

// EncodeColor converts a linear colour value to an sRGB encoded one.
func EncodeColor(c color.Color) color.RGBA64 {
	col, alpha := ColorFromLinearColor(c)
	return col.ToRGBA64(alpha)
}

// EncodeImage converts an image with linear colour into an sRGB encoded one.
//
// src is the linearised image to be encoded.
//
// dst is the image to write the result to, beginning at its origin.
//
// src and dst may be the same image.
//
// parallelism specifies the degree of parallel processing; a value of 4
// indicates that processing will be spread across four threads.
func EncodeImage(dst draw.Image, src image.Image, parallelism int) {
	linear.TransformImageColor(dst, src, parallelism, EncodeColor)
}

func encodedToLinear(v float32) float32 {
	if v <= 0.0031308*12.92 {
		return v / 12.92
	}
	return float32(math.Pow((float64(v)+0.055)/1.055, 2.4))
}

// LineariseColor converts an sRGB encoded colour into a linear one.
func LineariseColor(c color.Color) color.RGBA64 {
	col, alpha := ColorFromEncodedColor(c)
	return col.ToLinearRGBA64(alpha)
}

// LineariseImage converts an image with sRGB encoded colour to linear colour.
//
// src is the encoded image to be linearised.
//
// dst is the image to write the result to, beginning at its origin.
//
// src and dst may be the same image.
//
// parallelism specifies the degree of parallel processing; a value of 4
// indicates that processing will be spread across four threads.
func LineariseImage(dst draw.Image, src image.Image, parallelism int) {
	linear.TransformImageColor(dst, src, parallelism, LineariseColor)
}

func linearToEncoded(v float32) float32 {
	if v <= 0.0031308 {
		return v * 12.92
	}
	return float32(1.055*math.Pow(float64(v), 1/2.4) - 0.055)
}
