package vibrant

import (
	"image/color"
	"math"
)

// A Filter provides a mechanism for exercising fine-grained control over which colors
// are valid within a resulting Palette.
type Filter interface {
	// Hook to allow clients to be able filter colors from resulting palette.  Return true if the color is allowed, false if not.
	isAllowed(c color.Color) bool
}

// Default values for the default filter.
const (
	DefaultFilterMinLightness = 0.05
	DefaultFilterMaxLightness = 0.95
)

// DefaultFilter removes colors close to white and black.
var DefaultFilter = NewLightnessFilter(DefaultFilterMinLightness, DefaultFilterMaxLightness)

type defaultFilter struct {
	blackMaxLightness     float64
	whiteMinimumLightness float64
	quantizedBlackMask    uint16
	quantizedWhiteMask    uint16
}

func (f *defaultFilter) isAllowed(c color.Color) bool {
	// Short circuit for quantized colors to allow for faster black/white checking.
	if q, ok := c.(QuantizedColor); ok && !f.isAllowedQuantizedColor(q) {
		return false
	}
	hsl := HSLModel.Convert(c).(HSL)
	return !f.isWhite(hsl) && !f.isBlack(hsl)
}

func (f *defaultFilter) isAllowedQuantizedColor(q QuantizedColor) bool {
	return !(uint16(q)|f.quantizedBlackMask == f.quantizedBlackMask || uint16(q)&f.quantizedWhiteMask == f.quantizedWhiteMask)
}

func (f *defaultFilter) isWhite(hsl HSL) bool {
	return hsl.L >= f.whiteMinimumLightness
}

func (f *defaultFilter) isBlack(hsl HSL) bool {
	return hsl.L <= f.blackMaxLightness
}

// NewLightnessFilter creates a Filter that removes colors based on a lightness.
func NewLightnessFilter(blackMaxLightness float64, whiteMinimumLightness float64) Filter {
	quantizedBlackMask := uint16(uint8(math.Floor(float64(blackMaxLightness)*float64(255))) >> 3)
	quantizedWhiteMask := uint16(uint8(math.Ceil(float64(whiteMinimumLightness)*float64(255))) >> 3)
	return &defaultFilter{
		blackMaxLightness,
		whiteMinimumLightness,
		(quantizedBlackMask << 10) | (quantizedBlackMask << 5) | quantizedBlackMask,
		(quantizedWhiteMask << 10) | (quantizedWhiteMask << 5) | quantizedWhiteMask,
	}
}
