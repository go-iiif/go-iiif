package adobergb

import (
	"github.com/mandykoh/prism/ciexyy"
	"github.com/mandykoh/prism/linear"
	"image"
	"image/color"
	"image/draw"
	"math"
)

var PrimaryRed = ciexyy.Color{X: 0.64, Y: 0.33, YY: 1}
var PrimaryGreen = ciexyy.Color{X: 0.21, Y: 0.71, YY: 1}
var PrimaryBlue = ciexyy.Color{X: 0.15, Y: 0.06, YY: 1}
var StandardWhitePoint = ciexyy.D65

// EncodeColor converts a linear colour value to an Adobe RGB encoded one.
func EncodeColor(c color.Color) color.RGBA64 {
	col, alpha := ColorFromLinearColor(c)
	return col.ToRGBA64(alpha)
}

// EncodeImage converts an image with linear colour into an Adobe RGB encoded
// one.
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
	return float32(math.Pow(float64(v), 563.0/256))
}

// LineariseColor converts an Adobe RGB encoded colour into a linear one.
func LineariseColor(c color.Color) color.RGBA64 {
	col, alpha := ColorFromEncodedColor(c)
	return col.ToLinearRGBA64(alpha)
}

// LineariseImage converts an image with Adobe RGB encoded colour to linear
// colour.
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
	return float32(math.Pow(float64(v), 256.0/563))
}
