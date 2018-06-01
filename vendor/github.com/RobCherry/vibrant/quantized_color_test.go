package vibrant

import (
	"image/color"
	"strconv"
	"testing"
)

func TestNRGBAToQuantizedColorAndBack(t *testing.T) {
	tests := [][]color.Color{
		{color.NRGBA{0, 0, 0, 0xFF}, QuantizedColor(0), color.NRGBA{0, 0, 0, 0xFF}},
		{color.NRGBA{0x03, 0x03, 0x03, 0xFF}, QuantizedColor(0), color.NRGBA{0, 0, 0, 0xFF}},
		{color.NRGBA{0x04, 0x04, 0x04, 0xFF}, QuantizedColor(0x421), color.NRGBA{0x08, 0x08, 0x08, 0xFF}},
		{color.NRGBA{0x08, 0x08, 0x08, 0xFF}, QuantizedColor(0x421), color.NRGBA{0x08, 0x08, 0x08, 0xFF}},
		{color.NRGBA{0x0F, 0x0F, 0x0F, 0xFF}, QuantizedColor(0x421), color.NRGBA{0x08, 0x08, 0x08, 0xFF}},
		{color.NRGBA{0x10, 0x10, 0x10, 0xFF}, QuantizedColor(0x842), color.NRGBA{0x10, 0x10, 0x10, 0xFF}},
		{color.NRGBA{0xF0, 0xF0, 0xF0, 0xFF}, QuantizedColor(0x7BDE), color.NRGBA{0xF0, 0xF0, 0xF0, 0xFF}},
		{color.NRGBA{0xF3, 0xF3, 0xF3, 0xFF}, QuantizedColor(0x7BDE), color.NRGBA{0xF0, 0xF0, 0xF0, 0xFF}},
		{color.NRGBA{0xF4, 0xF4, 0xF4, 0xFF}, QuantizedColor(0x7FFF), color.NRGBA{0xF8, 0xF8, 0xF8, 0xFF}},
		{color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}, QuantizedColor(0x7FFF), color.NRGBA{0xF8, 0xF8, 0xF8, 0xFF}},
	}
	for _, test := range tests {
		originalValue := test[0]
		expectedQuantizedColor := test[1]
		expectedAproximateNRGBA := test[2]
		actualQuantizedColor := QuantizedColorModel.Convert(originalValue)
		if actualQuantizedColor != expectedQuantizedColor {
			t.Errorf("Color %v converted to %v instead of %v as expected.\n", originalValue, actualQuantizedColor, expectedQuantizedColor)
		} else if color.NRGBAModel.Convert(actualQuantizedColor) != expectedAproximateNRGBA {
			t.Errorf("Color %v converted to %v instead of %v as expected.\n", actualQuantizedColor, color.NRGBAModel.Convert(actualQuantizedColor), expectedAproximateNRGBA)
		}
	}
}

func TestQuantizedColorValue(t *testing.T) {
	tests := map[uint8]uint16{
		0: 0, 1: 0, 2: 0, 3: 0,
		4: 1, 5: 1, 6: 1, 7: 1,
		8: 1, 9: 1, 10: 1, 11: 1,
	}
	for value, expected := range tests {
		actual := quantizeColorValue(value)
		if actual != expected {
			t.Errorf("Value %08s converted to %05s instead of %05s as expected.\n", strconv.FormatUint(uint64(value), 2), strconv.FormatUint(uint64(actual), 2), strconv.FormatUint(uint64(expected), 2))
		}
	}
}

func TestQuantizedColor_SwapRedGreen(t *testing.T) {
	tests := [][]QuantizedColor{
		{QuantizedColor(0x7C00), QuantizedColor(0x3E0)},
		{QuantizedColor(0x7FE0), QuantizedColor(0x7FE0)},
		{QuantizedColor(0x1F), QuantizedColor(0x1F)},
	}
	for _, test := range tests {
		swappedValue := test[0].SwapRedGreen()
		unswappedValue := swappedValue.SwapRedGreen()
		expectedValue := test[1]
		if swappedValue != expectedValue {
			t.Errorf("Color %v converted to %v instead of %v as expected.\n", test[0], swappedValue, expectedValue)
		} else if unswappedValue != test[0] {
			t.Errorf("Color %v converted to %v instead of %v as expected.\n", swappedValue, unswappedValue, test[0])
		}
	}
}

func TestQuantizedColor_SwapRedBlue(t *testing.T) {
	tests := [][]QuantizedColor{
		{QuantizedColor(0x7C00), QuantizedColor(0x1F)},
		{QuantizedColor(0x7C1F), QuantizedColor(0x7C1F)},
		{QuantizedColor(0x3E0), QuantizedColor(0x3E0)},
	}
	for _, test := range tests {
		swappedValue := test[0].SwapRedBlue()
		unswappedValue := swappedValue.SwapRedBlue()
		expectedValue := test[1]
		if swappedValue != expectedValue {
			t.Errorf("Color %v converted to %v instead of %v as expected.\n", test[0], swappedValue, expectedValue)
		} else if unswappedValue != test[0] {
			t.Errorf("Color %v converted to %v instead of %v as expected.\n", swappedValue, unswappedValue, test[0])
		}
	}
}
