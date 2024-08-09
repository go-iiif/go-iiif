package linear

import (
	"image/color"
)

// RGB represents a linear normalised RGB colour value in an unspecified colour
// space.
type RGB struct {
	R float32
	G float32
	B float32
}

// Luminance returns the perceptual luminance of this colour.
func (c RGB) Luminance() float32 {
	return 0.2126*c.R + 0.7152*c.G + 0.0722*c.B
}

// ToEncodedNRGBA returns an encoded 8-bit NRGBA representation of this colour
// suitable for use with instances of image.NRGBA.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
//
// trcEncode is a tonal response curve encoding function.
func (c RGB) ToEncodedNRGBA(alpha float32, trcEncode func(float32) uint8) color.NRGBA {
	return color.NRGBA{
		R: trcEncode(c.R),
		G: trcEncode(c.G),
		B: trcEncode(c.B),
		A: NormalisedTo8Bit(alpha),
	}
}

// ToEncodedRGBA returns an encoded 8-bit RGBA representation of this colour
// suitable for use with instances of image.RGBA.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
//
// trcEncode is a tonal response curve encoding function.
func (c RGB) ToEncodedRGBA(alpha float32, trcEncode func(float32) uint8) color.RGBA {
	return color.RGBA{
		R: trcEncode(c.R * alpha),
		G: trcEncode(c.G * alpha),
		B: trcEncode(c.B * alpha),
		A: NormalisedTo8Bit(alpha),
	}
}

// ToEncodedRGBA64 returns an encoded 16-bit RGBA representation of this colour
// suitable for use with instances of image.RGBA64.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
//
// trcEncode is a tonal response curve encoding function.
func (c RGB) ToEncodedRGBA64(alpha float32, trcEncode func(float32) uint16) color.RGBA64 {
	return color.RGBA64{
		R: trcEncode(c.R * alpha),
		G: trcEncode(c.G * alpha),
		B: trcEncode(c.B * alpha),
		A: NormalisedTo16Bit(alpha),
	}
}

// ToLinearRGBA64 returns a linear 16-bit RGBA representation of this colour
// suitable for use with instances of image.RGBA64.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
func (c RGB) ToLinearRGBA64(alpha float32) color.RGBA64 {
	return color.RGBA64{
		R: NormalisedTo16Bit(c.R * alpha),
		G: NormalisedTo16Bit(c.G * alpha),
		B: NormalisedTo16Bit(c.B * alpha),
		A: NormalisedTo16Bit(alpha),
	}
}

// RGBFromEncoded returns a normalised RGB instance representing the specified
// color.Color value. The alpha component is returned as a normalised value in
// the range 0.0-1.0.
//
// c is assumed to be an encoded colour.
func RGBFromEncoded(c color.Color, trcDecode func(uint16) float32) (col RGB, alpha float32) {
	r, g, b, a := c.RGBA()

	if a == 0 {
		return RGB{}, 0
	}

	alpha = float32(a) / 65535

	return RGB{
			R: trcDecode(uint16(r)) / alpha,
			G: trcDecode(uint16(g)) / alpha,
			B: trcDecode(uint16(b)) / alpha,
		},
		alpha
}

// RGBFromLinear returns a normalised RGB instance representing the specified
// color.Color value. The alpha component is returned as a normalised value in
// the range 0.0-1.0.
//
// c is assumed to be a linear colour.
func RGBFromLinear(c color.Color) (col RGB, alpha float32) {
	r, g, b, a := c.RGBA()

	if a == 0 {
		return RGB{}, 0
	}

	alpha = float32(a)

	return RGB{
			R: float32(r) / alpha,
			G: float32(g) / alpha,
			B: float32(b) / alpha,
		},
		alpha / 65535
}
