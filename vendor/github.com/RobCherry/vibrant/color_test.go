package vibrant

import (
	"image/color"
	"testing"
)

func TestNRGBAToHSLAndBack(t *testing.T) {
	tests := [][]color.Color{
		{color.NRGBA{0, 0, 0, 255}, HSL{0.0, 0.0, 0.0, 255}},
		{color.NRGBA{63, 63, 63, 255}, HSL{0.0, 0.0, 0.24705882352941178, 255}},
		{color.NRGBA{127, 127, 127, 255}, HSL{0.0, 0.0, 0.4980392156862745, 255}},
		{color.NRGBA{255, 255, 255, 255}, HSL{0.0, 0.0, 1.0, 255}},
		{color.NRGBA{255, 0, 0, 255}, HSL{0.0, 1.0, 0.5, 255}},
		{color.NRGBA{0, 255, 0, 255}, HSL{120.0, 1.0, 0.5, 255}},
		{color.NRGBA{0, 0, 255, 255}, HSL{240.0, 1.0, 0.5, 255}},
	}
	for _, test := range tests {
		value := test[0]
		expected := test[1]
		actual := HSLModel.Convert(value)
		if actual != expected {
			t.Errorf("Color %v converted to %v instead of %v as expected.\n", value, actual, expected)
		} else if color.NRGBAModel.Convert(actual) != value {
			t.Errorf("Color %v did not convert back to %v as expected.\n", actual, value)
		}
	}
}

func TestNRGBAToRGBAIntAndBack(t *testing.T) {
	tests := [][]color.Color{
		{color.NRGBA{0, 0, 0, 255}, RGBAInt(0xFF000000)},
		{color.NRGBA{63, 63, 63, 255}, RGBAInt(0xFF3F3F3F)},
		{color.NRGBA{127, 127, 127, 255}, RGBAInt(0xFF7F7F7F)},
		{color.NRGBA{255, 255, 255, 255}, RGBAInt(0xFFFFFFFF)},
		{color.NRGBA{255, 0, 0, 255}, RGBAInt(0xFFFF0000)},
		{color.NRGBA{0, 255, 0, 255}, RGBAInt(0xFF00FF00)},
		{color.NRGBA{0, 0, 255, 255}, RGBAInt(0xFF0000FF)},
	}
	for _, test := range tests {
		value := test[0]
		expected := test[1]
		actual := RGBAIntModel.Convert(value)
		if actual != expected {
			t.Errorf("Color %v converted to %v instead of %v as expected.\n", value, actual, expected)
		} else if color.NRGBAModel.Convert(actual) != value {
			t.Errorf("Color %v did not convert back to %v as expected.\n", actual, value)
		}
	}
}
