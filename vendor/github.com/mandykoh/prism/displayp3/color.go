package displayp3

import (
	"github.com/mandykoh/prism/ciexyz"
	"github.com/mandykoh/prism/linear"
	"github.com/mandykoh/prism/srgb"
	"image/color"
)

// Color represents a linear normalised colour in Display P3 space.
type Color struct {
	linear.RGB
}

// ToNRGBA returns an encoded 8-bit NRGBA representation of this colour suitable
// for use with instances of image.NRGBA.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
func (c Color) ToNRGBA(alpha float32) color.NRGBA {
	return c.RGB.ToEncodedNRGBA(alpha, srgb.To8Bit)
}

// ToRGBA returns an encoded 8-bit RGBA representation of this colour suitable
// for use with instances of image.RGBA.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
func (c Color) ToRGBA(alpha float32) color.RGBA {
	return c.RGB.ToEncodedRGBA(alpha, srgb.To8Bit)
}

// ToRGBA64 returns an encoded 16-bit RGBA representation of this colour
// suitable for use with instances of image.RGBA64.
//
// alpha is the normalised alpha value and will be clipped to 0.0–1.0.
func (c Color) ToRGBA64(alpha float32) color.RGBA64 {
	return c.RGB.ToEncodedRGBA64(alpha, srgb.To16Bit)
}

// ToXYZ returns a CIE XYZ representation of this colour.
func (c Color) ToXYZ() ciexyz.Color {
	return ciexyz.Color{
		X: c.R*0.48656856264244125 + c.G*0.2656727168458704 + c.B*0.19818726598669462,
		Y: c.R*0.22897344124350177 + c.G*0.6917516599220641 + c.B*0.07927489883443435,
		Z: c.R*0 + c.G*0.04511425370425419 + c.B*1.0437861931875305,
	}
}

// ColorFromEncodedColor creates a Color instance from a Display P3 encoded
// color.Color value. The alpha value is returned as a normalised value between
// 0.0–1.0.
func ColorFromEncodedColor(c color.Color) (col Color, alpha float32) {
	rgb, a := linear.RGBFromEncoded(c, srgb.From16Bit)
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
// as Display P3 encoded. The alpha value is returned as a normalised
// value between 0.0–1.0.
func ColorFromNRGBA(c color.NRGBA) (col Color, alpha float32) {
	return Color{
		RGB: linear.RGB{
			R: srgb.From8Bit(c.R),
			G: srgb.From8Bit(c.G),
			B: srgb.From8Bit(c.B),
		},
	}, float32(c.A) / 255
}

// ColorFromRGBA creates a Color instance by interpreting an 8-bit RGBA colour
// as Display P3 encoded. The alpha value is returned as a normalised value
// between 0.0–1.0.
func ColorFromRGBA(c color.RGBA) (col Color, alpha float32) {
	if c.A == 0 {
		return Color{}, 0
	}

	alpha = float32(c.A) / 255

	return Color{
			RGB: linear.RGB{
				R: srgb.From8Bit(c.R) / alpha,
				G: srgb.From8Bit(c.G) / alpha,
				B: srgb.From8Bit(c.B) / alpha,
			},
		},
		alpha
}

// ColorFromXYZ creates a Display P3 Color instance from a CIE XYZ colour.
func ColorFromXYZ(c ciexyz.Color) Color {
	return ColorFromLinear(
		c.X*2.493509087331807+c.Y*-0.931388074532663+c.Z*-0.40271279318557973,
		c.X*-0.8294731994547587+c.Y*1.7626305488413623+c.Z*0.0236242511428412,
		c.X*0.03585127357050431+c.Y*-0.07618395633732165+c.Z*0.9570295296681479,
	)
}
