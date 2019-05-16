package vibrant

import (
	"image/color"
)

// Constants used for manipulating a QuantizedColor.
const (
	quantizeWordWidth = 5
	quantizeWordMask  = (1 << quantizeWordWidth) - 1
	shouldRoundUpMask = 1 << ((8 - quantizeWordWidth) - 1)
	roundUpMask       = shouldRoundUpMask << 1
)

// QuantizedColorSlice attaches the methods of sort.Interface to []QuantizedColor, sorting in increasing order.
type QuantizedColorSlice []QuantizedColor

func (s QuantizedColorSlice) Len() int           { return len(s) }
func (s QuantizedColorSlice) Less(i, j int) bool { return uint16(s[i]) < uint16(s[j]) }
func (s QuantizedColorSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// QuantizedColorGenerator creates a new QuantizedColor from a given red, green, and blue value.
var QuantizedColorGenerator = func(r, g, b uint8) QuantizedColor {
	quantizedRed := quantizeColorValue(r)
	quantizedGreen := quantizeColorValue(g)
	quantizedBlue := quantizeColorValue(b)
	return QuantizedColor((quantizedRed << (quantizeWordWidth + quantizeWordWidth)) | (quantizedGreen << quantizeWordWidth) | quantizedBlue)
}

// QuantizedColorModel is the color.Model for the QuantizedColor type.
var QuantizedColorModel = color.ModelFunc(func(c color.Color) color.Color {
	if _, ok := c.(QuantizedColor); ok {
		return c
	}
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	return QuantizedColorGenerator(nrgba.R, nrgba.G, nrgba.B)
})

// QuantizedColor represents a reduced RGB color space.
type QuantizedColor uint16

// RGBA implements the color.Color interface.
func (q QuantizedColor) RGBA() (uint32, uint32, uint32, uint32) {
	r := uint32(q.ApproximateRed())
	r |= r << 8
	g := uint32(q.ApproximateGreen())
	g |= g << 8
	b := uint32(q.ApproximateBlue())
	b |= b << 8
	a := uint32(0xFFFF)
	return r, g, b, a
}

// ApproximateRGBA is the approximate RGBA value of the quantized color.
func (q QuantizedColor) ApproximateRGBA() uint32 {
	r := uint32(q.ApproximateRed())
	g := uint32(q.ApproximateGreen())
	b := uint32(q.ApproximateBlue())
	a := uint32(0xFF)
	return (a << 24) | (r << 16) | (g << 8) | b
}

// QuantizedRed is the red component of the quantized color.
func (q QuantizedColor) QuantizedRed() uint8 {
	return uint8((q >> (quantizeWordWidth + quantizeWordWidth)) & quantizeWordMask)
}

// QuantizedGreen is the green component of a quantized color.
func (q QuantizedColor) QuantizedGreen() uint8 {
	return uint8((q >> quantizeWordWidth) & quantizeWordMask)
}

// QuantizedBlue is the blue component of a quantized color.
func (q QuantizedColor) QuantizedBlue() uint8 {
	return uint8(q & quantizeWordMask)
}

// ApproximateRed is the approximate red component of the quantized color.
func (q QuantizedColor) ApproximateRed() uint8 {
	return modifyWordWidth(q.QuantizedRed(), quantizeWordWidth, 8)
}

// ApproximateGreen is the approximate green component of a quantized color.
func (q QuantizedColor) ApproximateGreen() uint8 {
	return modifyWordWidth(q.QuantizedGreen(), quantizeWordWidth, 8)
}

// ApproximateBlue is the approximate blue component of a quantized color.
func (q QuantizedColor) ApproximateBlue() uint8 {
	return modifyWordWidth(q.QuantizedBlue(), quantizeWordWidth, 8)
}

// SwapRedGreen returns a new QuantizedColor whose red and green color components have been swapped.
func (q QuantizedColor) SwapRedGreen() QuantizedColor {
	return QuantizedColor(uint16(q.QuantizedGreen())<<(quantizeWordWidth+quantizeWordWidth) | uint16(q.QuantizedRed())<<quantizeWordWidth | uint16(q.QuantizedBlue()))
}

// SwapRedBlue returns a new QuantizedColor whose red and blue color components have been swapped.
func (q QuantizedColor) SwapRedBlue() QuantizedColor {
	return QuantizedColor(uint16(q.QuantizedBlue())<<(quantizeWordWidth+quantizeWordWidth) | uint16(q.QuantizedGreen())<<quantizeWordWidth | uint16(q.QuantizedRed()))
}

func quantizeColorValue(value uint8) uint16 {
	if value&shouldRoundUpMask == shouldRoundUpMask {
		value = value | roundUpMask
	}
	return uint16(modifyWordWidth(value, 8, quantizeWordWidth))
}

func modifyWordWidth(value uint8, currentWidth uint8, targetWidth uint8) uint8 {
	var modifiedValue uint8
	if targetWidth > currentWidth {
		// If we're approximating up in word width, we'll shift up
		modifiedValue = value << (targetWidth - currentWidth)
	} else {
		// Else, we will just shift and keep the MSB
		modifiedValue = value >> (currentWidth - targetWidth)
	}
	return modifiedValue & ((1 << targetWidth) - 1)
}
