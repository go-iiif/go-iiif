package lut

import "github.com/mandykoh/prism/linear"

func BuildLinearTo8Bit(encode func(float32) float32) [512]uint8 {
	to8BitLUT := [512]uint8{}
	for i := range to8BitLUT {
		to8BitLUT[i] = linear.NormalisedTo8Bit(encode(float32(i) / 511))
	}
	return to8BitLUT
}

func BuildLinearTo16Bit(encode func(float32) float32) [65536]uint16 {
	to16BitLUT := [65536]uint16{}
	for i := range to16BitLUT {
		to16BitLUT[i] = linear.NormalisedTo16Bit(encode(float32(i) / 65535))
	}
	return to16BitLUT
}

func Build8BitToLinear(linearise func(float32) float32) [256]float32 {
	from8BitLUT := [256]float32{}
	for i := range from8BitLUT {
		from8BitLUT[i] = linearise(float32(i) / 255)
	}
	return from8BitLUT
}

func Build16BitToLinear(linearise func(float32) float32) [65536]float32 {
	from16BitLUT := [65536]float32{}
	for i := range from16BitLUT {
		from16BitLUT[i] = linearise(float32(i) / 65535)
	}
	return from16BitLUT
}
