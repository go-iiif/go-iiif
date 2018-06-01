package vibrant

import (
	"testing"
)

// A vBox wil only one color can not be split.
func TestVBox_Split1(t *testing.T) {
	tests := [][]QuantizedColor{
		{
			QuantizedColorGenerator(1<<3, 0, 0),
			QuantizedColorGenerator(1<<3, 0, 0),
		},
		{
			QuantizedColorGenerator(0, 1<<3, 0),
			QuantizedColorGenerator(0, 1<<3, 0),
		},
		{
			QuantizedColorGenerator(0, 0, 1<<3),
			QuantizedColorGenerator(0, 0, 1<<3),
		},
	}
	for _, colors := range tests {
		box := newVBox(colors, generateHistorgram(colors), 0, uint32(len(colors))-1)
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 1, box.upperIndex)
		refuteCanSplit(t, box)
	}
}

// Order of the colors in the VBox does not matter because they are sorted before the box is split.
func TestVBox_Split2(t *testing.T) {
	tests := [][]QuantizedColor{
		{
			QuantizedColorGenerator(1<<3, 0, 0),
			QuantizedColorGenerator(2<<3, 0, 0),
			QuantizedColorGenerator(3<<3, 0, 0),
			QuantizedColorGenerator(4<<3, 0, 0),
		},
		{
			QuantizedColorGenerator(4<<3, 0, 0),
			QuantizedColorGenerator(3<<3, 0, 0),
			QuantizedColorGenerator(2<<3, 0, 0),
			QuantizedColorGenerator(1<<3, 0, 0),
		},
		{
			QuantizedColorGenerator(1<<3, 0, 0),
			QuantizedColorGenerator(3<<3, 0, 0),
			QuantizedColorGenerator(2<<3, 0, 0),
			QuantizedColorGenerator(4<<3, 0, 0),
		},
		{
			QuantizedColorGenerator(4<<3, 0, 0),
			QuantizedColorGenerator(2<<3, 0, 0),
			QuantizedColorGenerator(3<<3, 0, 0),
			QuantizedColorGenerator(1<<3, 0, 0),
		},
		{
			QuantizedColorGenerator(0, 0, 4<<3),
			QuantizedColorGenerator(0, 0, 2<<3),
			QuantizedColorGenerator(0, 0, 3<<3),
			QuantizedColorGenerator(0, 0, 1<<3),
		},
		{
			QuantizedColorGenerator(0, 1<<3, 0),
			QuantizedColorGenerator(0, 2<<3, 0),
			QuantizedColorGenerator(0, 3<<3, 0),
			QuantizedColorGenerator(0, 4<<3, 0),
		},
		{
			QuantizedColorGenerator(0, 4<<3, 0),
			QuantizedColorGenerator(0, 3<<3, 0),
			QuantizedColorGenerator(0, 2<<3, 0),
			QuantizedColorGenerator(0, 1<<3, 0),
		},
		{
			QuantizedColorGenerator(0, 1<<3, 0),
			QuantizedColorGenerator(0, 3<<3, 0),
			QuantizedColorGenerator(0, 2<<3, 0),
			QuantizedColorGenerator(0, 4<<3, 0),
		},
		{
			QuantizedColorGenerator(0, 4<<3, 0),
			QuantizedColorGenerator(0, 2<<3, 0),
			QuantizedColorGenerator(0, 3<<3, 0),
			QuantizedColorGenerator(0, 1<<3, 0),
		},
		{
			QuantizedColorGenerator(0, 0, 1<<3),
			QuantizedColorGenerator(0, 0, 2<<3),
			QuantizedColorGenerator(0, 0, 3<<3),
			QuantizedColorGenerator(0, 0, 4<<3),
		},
		{
			QuantizedColorGenerator(0, 0, 4<<3),
			QuantizedColorGenerator(0, 0, 3<<3),
			QuantizedColorGenerator(0, 0, 2<<3),
			QuantizedColorGenerator(0, 0, 1<<3),
		},
		{
			QuantizedColorGenerator(0, 0, 1<<3),
			QuantizedColorGenerator(0, 0, 3<<3),
			QuantizedColorGenerator(0, 0, 2<<3),
			QuantizedColorGenerator(0, 0, 4<<3),
		},
	}
	for _, colors := range tests {
		box := newVBox(colors, generateHistorgram(colors), 0, uint32(len(colors))-1)
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 3, box.upperIndex)
		assertCanSplit(t, box)
		splitBox, err := box.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 1, box.upperIndex)
		assertIndexMatches(t, "lower", 2, splitBox.lowerIndex)
		assertIndexMatches(t, "upper", 3, splitBox.upperIndex)
		assertCanSplit(t, box)
		assertCanSplit(t, splitBox)
		splitBox2, err := box.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 0, box.upperIndex)
		assertIndexMatches(t, "lower", 1, splitBox2.lowerIndex)
		assertIndexMatches(t, "upper", 1, splitBox2.upperIndex)
		refuteCanSplit(t, box)
		refuteCanSplit(t, splitBox2)
		splitBox3, err := splitBox.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 2, splitBox.lowerIndex)
		assertIndexMatches(t, "upper", 2, splitBox.upperIndex)
		assertIndexMatches(t, "lower", 3, splitBox3.lowerIndex)
		assertIndexMatches(t, "upper", 3, splitBox3.upperIndex)
		refuteCanSplit(t, splitBox)
		refuteCanSplit(t, splitBox3)
	}
}

func TestVBox_Split3(t *testing.T) {
	tests := [][]QuantizedColor{
		{
			QuantizedColorGenerator(1<<3, 0, 0),
			QuantizedColorGenerator(1<<3, 0, 0),
			QuantizedColorGenerator(2<<3, 0, 0),
			QuantizedColorGenerator(3<<3, 0, 0),
		},
		{
			QuantizedColorGenerator(0, 1<<3, 0),
			QuantizedColorGenerator(0, 1<<3, 0),
			QuantizedColorGenerator(0, 2<<3, 0),
			QuantizedColorGenerator(0, 3<<3, 0),
		},
		{
			QuantizedColorGenerator(0, 0, 1<<3),
			QuantizedColorGenerator(0, 0, 1<<3),
			QuantizedColorGenerator(0, 0, 2<<3),
			QuantizedColorGenerator(0, 0, 3<<3),
		},
	}
	for _, colors := range tests {
		box := newVBox(colors, generateHistorgram(colors), 0, uint32(len(colors))-1)
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 3, box.upperIndex)
		assertCanSplit(t, box)
		splitBox, err := box.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 1, box.upperIndex)
		assertIndexMatches(t, "lower", 2, splitBox.lowerIndex)
		assertIndexMatches(t, "upper", 3, splitBox.upperIndex)
		refuteCanSplit(t, box)
		assertCanSplit(t, splitBox)
		splitBox2, err := splitBox.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 2, splitBox.lowerIndex)
		assertIndexMatches(t, "upper", 2, splitBox.upperIndex)
		assertIndexMatches(t, "lower", 3, splitBox2.lowerIndex)
		assertIndexMatches(t, "upper", 3, splitBox2.upperIndex)
		refuteCanSplit(t, splitBox)
		refuteCanSplit(t, splitBox2)
	}
}

func TestVBox_Split4(t *testing.T) {
	tests := [][]QuantizedColor{
		{
			QuantizedColorGenerator(1<<3, 0, 0),
			QuantizedColorGenerator(2<<3, 0, 0),
			QuantizedColorGenerator(2<<3, 0, 0),
			QuantizedColorGenerator(3<<3, 0, 0),
		},
		{
			QuantizedColorGenerator(0, 1<<3, 0),
			QuantizedColorGenerator(0, 2<<3, 0),
			QuantizedColorGenerator(0, 2<<3, 0),
			QuantizedColorGenerator(0, 3<<3, 0),
		},
		{
			QuantizedColorGenerator(0, 0, 1<<3),
			QuantizedColorGenerator(0, 0, 2<<3),
			QuantizedColorGenerator(0, 0, 2<<3),
			QuantizedColorGenerator(0, 0, 3<<3),
		},
	}
	for _, colors := range tests {
		box := newVBox(colors, generateHistorgram(colors), 0, uint32(len(colors))-1)
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 3, box.upperIndex)
		assertCanSplit(t, box)
		splitBox, err := box.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 2, box.upperIndex)
		assertIndexMatches(t, "lower", 3, splitBox.lowerIndex)
		assertIndexMatches(t, "upper", 3, splitBox.upperIndex)
		assertCanSplit(t, box)
		refuteCanSplit(t, splitBox)
		splitBox2, err := box.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 0, box.upperIndex)
		assertIndexMatches(t, "lower", 1, splitBox2.lowerIndex)
		assertIndexMatches(t, "upper", 2, splitBox2.upperIndex)
		refuteCanSplit(t, box)
		refuteCanSplit(t, splitBox2)
	}
}

func TestVBox_Split5(t *testing.T) {
	tests := [][]QuantizedColor{
		{
			QuantizedColorGenerator(1<<3, 0, 0),
			QuantizedColorGenerator(2<<3, 0, 0),
			QuantizedColorGenerator(3<<3, 0, 0),
			QuantizedColorGenerator(3<<3, 0, 0),
		},
		{
			QuantizedColorGenerator(0, 1<<3, 0),
			QuantizedColorGenerator(0, 2<<3, 0),
			QuantizedColorGenerator(0, 3<<3, 0),
			QuantizedColorGenerator(0, 3<<3, 0),
		},
		{
			QuantizedColorGenerator(0, 0, 1<<3),
			QuantizedColorGenerator(0, 0, 2<<3),
			QuantizedColorGenerator(0, 0, 3<<3),
			QuantizedColorGenerator(0, 0, 3<<3),
		},
	}
	for _, colors := range tests {
		box := newVBox(colors, generateHistorgram(colors), 0, uint32(len(colors))-1)
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 3, box.upperIndex)
		assertCanSplit(t, box)
		splitBox, err := box.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 1, box.upperIndex)
		assertIndexMatches(t, "lower", 2, splitBox.lowerIndex)
		assertIndexMatches(t, "upper", 3, splitBox.upperIndex)
		assertCanSplit(t, box)
		refuteCanSplit(t, splitBox)
		splitBox2, err := box.Split()
		if err != nil {
			t.Fatal(err)
		}
		assertIndexMatches(t, "lower", 0, box.lowerIndex)
		assertIndexMatches(t, "upper", 0, box.upperIndex)
		assertIndexMatches(t, "lower", 1, splitBox2.lowerIndex)
		assertIndexMatches(t, "upper", 1, splitBox2.upperIndex)
		refuteCanSplit(t, box)
		refuteCanSplit(t, splitBox2)
	}
}

func TestVBox_Volume(t *testing.T) {
	tests := map[uint32][][]QuantizedColor{
		1: {
			{
				QuantizedColorGenerator(0, 0, 1<<3),
			},
			{
				QuantizedColorGenerator(0, 0, 1<<3),
				QuantizedColorGenerator(0, 0, 1<<3),
			},
			{
				QuantizedColorGenerator(1<<3, 2<<3, 3<<3),
			},
		},
		2: {
			{
				QuantizedColorGenerator(0, 0, 0),
				QuantizedColorGenerator(0, 0, 1<<3),
			},
			{
				QuantizedColorGenerator(0, 0, 1<<3),
				QuantizedColorGenerator(1<<3, 0, 1<<3),
			},
		},
		4: {
			{
				QuantizedColorGenerator(1<<3, 0, 0),
				QuantizedColorGenerator(0, 0, 1<<3),
			},
		},
		8: {
			{
				QuantizedColorGenerator(1<<3, 0, 0),
				QuantizedColorGenerator(0, 1<<3, 0),
				QuantizedColorGenerator(0, 0, 1<<3),
			},
		},
	}
	for expected, arrayOfColors := range tests {
		for _, colors := range arrayOfColors {
			box := newVBox(colors, generateHistorgram(colors), 0, uint32(len(colors))-1)
			if expected != box.Volume() {
				t.Fatalf("Expected volume %v != %v", expected, box.Volume())
			}
		}
	}
}

func TestVBox_Swatch(t *testing.T) {
	tests := map[RGBAInt][][]QuantizedColor{
		// vBox colors are quantized, so we need to adjust the expected values to account for it.
		RGBAInt(0xFFF80000): {
			{
				QuantizedColorGenerator(255, 0, 0),
			},
			{
				QuantizedColorGenerator(255, 0, 0),
				QuantizedColorGenerator(255, 0, 0),
			},
		},
		RGBAInt(0xFF00F800): {
			{
				QuantizedColorGenerator(0, 255, 0),
			},
			{
				QuantizedColorGenerator(0, 255, 0),
				QuantizedColorGenerator(0, 255, 0),
			},
		},
		RGBAInt(0xFF0000F8): {
			{
				QuantizedColorGenerator(0, 0, 255),
			},
			{
				QuantizedColorGenerator(0, 0, 255),
				QuantizedColorGenerator(0, 0, 255),
			},
		},
		RGBAInt(0xFFF8F400): {
			{
				QuantizedColorGenerator(0xFF, 0xFF, 0),
				QuantizedColorGenerator(0xFF, 0xF0, 0),
			},
			{
				QuantizedColorGenerator(0xFF, 0xF8, 0),
				QuantizedColorGenerator(0xFF, 0xF0, 0),
			},
		},
		RGBAInt(0xFFF2F600): {
			{
				QuantizedColorGenerator(0xF0, 0xFF, 0),
				QuantizedColorGenerator(0xF0, 0xFF, 0),
				QuantizedColorGenerator(0xFF, 0xF0, 0),
			},
		},
		RGBAInt(0xFFF1F700): {
			{
				QuantizedColorGenerator(0xF0, 0xFF, 0),
				QuantizedColorGenerator(0xF0, 0xFF, 0),
				QuantizedColorGenerator(0xF0, 0xFF, 0),
				QuantizedColorGenerator(0xFF, 0xF0, 0),
			},
		},
	}
	for expected, arrayOfColors := range tests {
		for _, colors := range arrayOfColors {
			box := newVBox(colors, generateHistorgram(colors), 0, uint32(len(colors))-1)
			swatch := box.Swatch()
			if expected != swatch.RGBAInt() {
				t.Errorf("Expected swatch %v != %v", expected, swatch.RGBAInt())
			}
		}
	}
}

func generateHistorgram(colors []QuantizedColor) []uint32 {
	histogram := make([]uint32, histogramSize)
	for _, color := range colors {
		histogram[color]++
	}
	return histogram
}

func assertIndexMatches(t *testing.T, upperOrLower string, expected, actual uint32) {
	if expected != actual {
		t.Errorf("Expected %s index %d != %d", upperOrLower, expected, actual)
	}
}

func assertCanSplit(t *testing.T, box *vBox) {
	if !box.CanSplit() {
		t.Fatal("Expected to be able to split box")
	}
}

func refuteCanSplit(t *testing.T, box *vBox) {
	if box.CanSplit() {
		t.Fatal("Expected to not be able to split box")
	}
}
