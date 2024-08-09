package srgb

import (
	"github.com/mandykoh/prism/ciexyz"
	"github.com/mandykoh/prism/linear"
	"image/color"
)

// Color represents a linear normalised colour in sRGB space.
type Color struct {
	linear.RGB
}

// ToNRGBA returns an encoded 8-bit NRGBA representation of this colour suitable
// for use with instances of image.NRGBA.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
func (c Color) ToNRGBA(alpha float32) color.NRGBA {
	return c.RGB.ToEncodedNRGBA(alpha, To8Bit)
}

// ToRGBA returns an encoded 8-bit RGBA representation of this colour suitable
// for use with instances of image.RGBA.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
func (c Color) ToRGBA(alpha float32) color.RGBA {
	return c.RGB.ToEncodedRGBA(alpha, To8Bit)
}

// ToRGBA64 returns an encoded 16-bit RGBA representation of this colour
// suitable for use with instances of image.RGBA64.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
func (c Color) ToRGBA64(alpha float32) color.RGBA64 {
	return c.RGB.ToEncodedRGBA64(alpha, To16Bit)
}

// ToXYZ returns a CIE XYZ representation of this colour.
func (c Color) ToXYZ() ciexyz.Color {
	return ciexyz.Color{
		X: c.R*0.41238652374145657 + c.G*0.35759149384555927 + c.B*0.18045052788799043,
		Y: c.R*0.21263680803732615 + c.G*0.7151829876911185 + c.B*0.07218020427155547,
		Z: c.R*0.019330619488581533 + c.G*0.11919711488225383 + c.B*0.9503727125209493,
	}
}

// ColorFromEncodedColor creates a Color instance from an sRGB encoded
// color.Color value. The alpha value is returned as a normalised value between
// 0.0–1.0.
func ColorFromEncodedColor(c color.Color) (col Color, alpha float32) {
	rgb, a := linear.RGBFromEncoded(c, From16Bit)
	return Color{rgb}, a
}

// ColorFromLinear creates a Color instance from a linear normalised RGB
// triplet.
func ColorFromLinear(r, g, b float32) Color {
	return Color{linear.RGB{R: r, G: g, B: b}}
}

// ColorFromLinearColor creates a Color instance from a linear color.Color
// value. The alpha value is returned as a normalised value between 0.0–1.0.
func ColorFromLinearColor(c color.Color) (col Color, alpha float32) {
	rgb, a := linear.RGBFromLinear(c)
	return Color{rgb}, a
}

// ColorFromNRGBA creates a Color instance by interpreting an 8-bit NRGBA colour
// as sRGB encoded. The alpha value is returned as a normalised value between
// 0.0–1.0.
func ColorFromNRGBA(c color.NRGBA) (col Color, alpha float32) {
	return Color{
		RGB: linear.RGB{
			R: From8Bit(c.R),
			G: From8Bit(c.G),
			B: From8Bit(c.B),
		},
	}, float32(c.A) / 255
}

// ColorFromRGBA creates a Color instance by interpreting an 8-bit RGBA colour
// as sRGB encoded. The alpha value is returned as a normalised value between
// 0.0–1.0.
func ColorFromRGBA(c color.RGBA) (col Color, alpha float32) {
	if c.A == 0 {
		return Color{}, 0
	}

	alpha = float32(c.A) / 255

	return Color{
			RGB: linear.RGB{
				R: From8Bit(c.R) / alpha,
				G: From8Bit(c.G) / alpha,
				B: From8Bit(c.B) / alpha,
			},
		},
		alpha
}

// ColorFromXYZ creates an sRGB Color instance from a CIE XYZ colour.
func ColorFromXYZ(c ciexyz.Color) Color {
	return ColorFromLinear(
		c.X*3.241003600540255+c.Y*-1.5373991710891957+c.Z*-0.49861598312439226,
		c.X*-0.9692242864995344+c.Y*1.8759299885141119+c.Z*0.04155424903337176,
		c.X*0.05563936186796137+c.Y*-0.20401108051523723+c.Z*1.0571488385644063,
	)
}
