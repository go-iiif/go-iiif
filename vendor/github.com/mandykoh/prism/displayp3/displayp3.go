package displayp3

import (
	"github.com/mandykoh/prism/ciexyy"
	"github.com/mandykoh/prism/linear"
	"image"
	"image/color"
	"image/draw"
)

var PrimaryRed = ciexyy.Color{X: 0.68, Y: 0.32, YY: 1}
var PrimaryGreen = ciexyy.Color{X: 0.265, Y: 0.69, YY: 1}
var PrimaryBlue = ciexyy.Color{X: 0.15, Y: 0.06, YY: 1}
var StandardWhitePoint = ciexyy.D65

// EncodeColor converts a linear colour value to a Display P3 encoded one.
func EncodeColor(c color.Color) color.RGBA64 {
	col, alpha := ColorFromLinearColor(c)
	return col.ToRGBA64(alpha)
}

// EncodeImage converts an image with linear colour into a Display P3 encoded
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

// LineariseColor converts a Display P3 encoded colour into a linear one.
func LineariseColor(c color.Color) color.RGBA64 {
	col, alpha := ColorFromEncodedColor(c)
	return col.ToLinearRGBA64(alpha)
}

// LineariseImage converts an image with Display P3 encoded colour to linear
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
