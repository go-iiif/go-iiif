package vibrant

import (
	"image/color"
	"strconv"
	"testing"
)

func TestDefaultFilter_isAllowed(t *testing.T) {
	testData := map[color.Color]bool{
		HSL{0, 1.0, 0.049, 255}: false,
		HSL{0, 1.0, 0.050, 255}: false,
		HSL{0, 1.0, 0.051, 255}: true,
		HSL{0, 1.0, 0.949, 255}: true,
		HSL{0, 1.0, 0.950, 255}: false,
		HSL{0, 1.0, 0.951, 255}: false,
	}
	for c, expected := range testData {
		if expected != DefaultFilter.isAllowed(c) {
			if expected {
				t.Errorf("Expected color %v to be allowed.\n", c)
			} else {
				t.Errorf("Expected color %v to be filtered.\n", c)
			}
		}
	}
}

func TestDefaultFilter_isAllowedQuantizedColor(t *testing.T) {
	testData := map[QuantizedColor]bool{
		QuantizedColorGenerator(0, 0, 0):       false,
		QuantizedColorGenerator(15, 15, 15):    false,
		QuantizedColorGenerator(15, 15, 16):    true,
		QuantizedColorGenerator(15, 16, 16):    true,
		QuantizedColorGenerator(16, 16, 16):    true,
		QuantizedColorGenerator(239, 239, 239): true,
		QuantizedColorGenerator(239, 239, 240): true,
		QuantizedColorGenerator(239, 240, 240): true,
		QuantizedColorGenerator(240, 240, 240): false,
		QuantizedColorGenerator(255, 255, 255): false,
	}
	filter := DefaultFilter.(*defaultFilter)
	for q, expected := range testData {
		if expected != filter.isAllowedQuantizedColor(q) {
			t.Logf("Black mask: %016s", strconv.FormatInt(int64(filter.quantizedBlackMask), 2))
			t.Logf("White mask: %016s", strconv.FormatInt(int64(filter.quantizedWhiteMask), 2))
			if expected {
				t.Errorf("Expected quantized color %v (%016s / 0x%x) to be allowed.\n", q, strconv.FormatInt(int64(q), 2), q)
			} else {
				t.Errorf("Expected quantized color %v (%016s / 0x%x) to be filtered.\n", q, strconv.FormatInt(int64(q), 2), q)
			}
		}
	}
}
