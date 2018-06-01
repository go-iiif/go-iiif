package vibrant

import (
	"image/color"
)

// Swatch represents a color swatch generated from an image's palette.
type Swatch struct {
	color      color.Color
	rgbaInt    *RGBAInt
	hsl        *HSL
	population uint32
}

// NewSwatch creates a new Swatch from a color and population.
func NewSwatch(color color.Color, population uint32) *Swatch {
	return &Swatch{
		color:      color,
		population: population,
	}
}

// RGBAInt returns this swatch's color value as an RGBAInt.
func (s *Swatch) RGBAInt() RGBAInt {
	if s.rgbaInt == nil {
		rgbaInt := RGBAIntModel.Convert(s.color).(RGBAInt)
		s.rgbaInt = &rgbaInt
	}
	return *s.rgbaInt
}

// HSL returns this swatch's color value as an HSL.
func (s *Swatch) HSL() HSL {
	if s.hsl == nil {
		hsl := HSLModel.Convert(s.color).(HSL)
		s.hsl = &hsl
	}
	return *s.hsl
}

// Color returns this swatch's color value.
func (s Swatch) Color() color.Color {
	return s.color
}

// Population returns the number of pixels represented by this swatch.
func (s Swatch) Population() uint32 {
	return s.population
}
