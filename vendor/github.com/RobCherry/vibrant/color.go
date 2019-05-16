package vibrant

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// RGBAInt represents a packed RGBA color.
type RGBAInt uint32

// RGBA implements the color.Color interface.
func (c RGBAInt) RGBA() (uint32, uint32, uint32, uint32) {
	alphaComponent := uint32(c) >> 24
	r := (uint32(c) >> 16) & 0xFF
	r |= r << 8
	r *= alphaComponent
	r /= 0xFF
	g := (uint32(c) >> 8) & 0xFF
	g |= g << 8
	g *= alphaComponent
	g /= 0xFF
	b := uint32(c) & 0xFF
	b |= b << 8
	b *= alphaComponent
	b /= 0xFF
	a := alphaComponent
	a |= a << 8
	return r, g, b, a
}

// PackedRGBA is the packed int representing the RGBA value.
func (c RGBAInt) PackedRGBA() uint32 {
	return uint32(c)
}

// PackedRGB is the packed int representing the RGB value (ignores the alpha channel).
func (c RGBAInt) PackedRGB() uint32 {
	return c.PackedRGBA() | (uint32(0xFF) << 24)
}

func (c RGBAInt) String() string {
	return fmt.Sprintf("0x%06s", strings.ToUpper(strconv.FormatUint(uint64(c.PackedRGBA()), 16)))
}

// RGBAIntModel is the color.Model for the RGBAInt type.
var RGBAIntModel = color.ModelFunc(func(c color.Color) color.Color {
	if _, ok := c.(RGBAInt); ok {
		return c
	}
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	return RGBAInt((uint32(nrgba.A) << 24) | (uint32(nrgba.R) << 16) | (uint32(nrgba.G) << 8) | uint32(nrgba.B))
})

// HSL represents the HSL value for an RGBA color.
type HSL struct {
	H, S, L float64
	A       uint8
}

// RGBA implements the color.Color interface.
func (c HSL) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := hslToRGB(c.H, c.S, c.L)
	return color.NRGBA{r, g, b, c.A}.RGBA()
}

func (c HSL) String() string {
	return fmt.Sprintf("HSL %0.0f %0.2f %0.2f", c.H, c.S, c.L)
}

// HSLModel is the color.Model for the HSL type.
var HSLModel = color.ModelFunc(func(c color.Color) color.Color {
	if _, ok := c.(HSL); ok {
		return c
	}
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	h, s, l := rgbToHSL(nrgba.R, nrgba.G, nrgba.B)
	return HSL{h, s, l, nrgba.A}
})

// Returns the Hue [0..360], Saturation and Lightness [0..1] of the color.
func rgbToHSL(r, g, b uint8) (h, s, l float64) {
	fr := float64(r) / 255.0
	fg := float64(g) / 255.0
	fb := float64(b) / 255.0

	min := math.Min(math.Min(fr, fg), fb)
	max := math.Max(math.Max(fr, fg), fb)

	l = (max + min) / 2
	if min == max {
		s = 0
		h = 0
	} else {
		if l < 0.5 {
			s = (max - min) / (max + min)
		} else {
			s = (max - min) / (2.0 - max - min)
		}
		if max == fr {
			h = (fg - fb) / (max - min)
		} else if max == fg {
			h = 2.0 + (fb-fr)/(max-min)
		} else {
			h = 4.0 + (fr-fg)/(max-min)
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	return h, s, l
}

// Returns the RGB [0..255] values given a Hue [0..360], Saturation and Lightness [0..1]
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	if s == 0 {
		return clampUint8(uint8(roundFloat64(l*255.0)), 0, 255), clampUint8(uint8(roundFloat64(l*255.0)), 0, 255), clampUint8(uint8(roundFloat64(l*255.0)), 0, 255)
	}

	var (
		r, g, b    float64
		t1, t2     float64
		tr, tg, tb float64
	)

	if l < 0.5 {
		t1 = l * (1.0 + s)
	} else {
		t1 = l + s - l*s
	}

	t2 = 2*l - t1
	h = h / 360
	tr = h + 1.0/3.0
	tg = h
	tb = h - 1.0/3.0

	if tr < 0 {
		tr++
	}
	if tr > 1 {
		tr--
	}
	if tg < 0 {
		tg++
	}
	if tg > 1 {
		tg--
	}
	if tb < 0 {
		tb++
	}
	if tb > 1 {
		tb--
	}

	// Red
	if 6*tr < 1 {
		r = t2 + (t1-t2)*6*tr
	} else if 2*tr < 1 {
		r = t1
	} else if 3*tr < 2 {
		r = t2 + (t1-t2)*(2.0/3.0-tr)*6
	} else {
		r = t2
	}

	// Green
	if 6*tg < 1 {
		g = t2 + (t1-t2)*6*tg
	} else if 2*tg < 1 {
		g = t1
	} else if 3*tg < 2 {
		g = t2 + (t1-t2)*(2.0/3.0-tg)*6
	} else {
		g = t2
	}

	// Blue
	if 6*tb < 1 {
		b = t2 + (t1-t2)*6*tb
	} else if 2*tb < 1 {
		b = t1
	} else if 3*tb < 2 {
		b = t2 + (t1-t2)*(2.0/3.0-tb)*6
	} else {
		b = t2
	}

	return clampUint8(uint8(roundFloat64(r*255.0)), 0, 255), clampUint8(uint8(roundFloat64(g*255.0)), 0, 255), clampUint8(uint8(roundFloat64(b*255.0)), 0, 255)
}
