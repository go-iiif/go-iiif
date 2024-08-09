package adobergb

import (
	"github.com/mandykoh/prism/linear"
	"github.com/mandykoh/prism/linear/lut"
	"sync"
)

var init8BitLUTsOnce sync.Once
var linearToEncoded8LUT []uint8
var encoded8ToLinearLUT []float32

var initTo16BitLUTOnce sync.Once
var linearToEncoded16LUT []uint16

var initFrom16BitLUTOnce sync.Once
var encoded16ToLinearLUT []float32

func init() {
	init8BitLUTsOnce.Do(func() {
		to8BitLUT := lut.BuildLinearTo8Bit(linearToEncoded)
		linearToEncoded8LUT = to8BitLUT[:]

		from8BitLUT := lut.Build8BitToLinear(encodedToLinear)
		encoded8ToLinearLUT = from8BitLUT[:]
	})
}

// From8Bit converts an 8-bit Adobe RGB encoded value to a normalised linear
// value between 0.0 and 1.0.
//
// This implementation uses a fast look-up table without sacrificing accuracy.
func From8Bit(v uint8) float32 {
	return encoded8ToLinearLUT[v]
}

// From16Bit converts a 16-bit Adobe RGB encoded value to a normalised linear
// value between 0.0 and 1.0.
//
// This implementation uses a fast look-up table without sacrificing accuracy.
func From16Bit(v uint16) float32 {
	if encoded16ToLinearLUT != nil {
		return encoded16ToLinearLUT[v]
	}
	return from16BitAndInitLUT(v)
}

func from16BitAndInitLUT(v uint16) float32 {
	initFrom16BitLUTOnce.Do(func() {
		from16BitLUT := lut.Build16BitToLinear(encodedToLinear)
		encoded16ToLinearLUT = from16BitLUT[:]
	})
	return encoded16ToLinearLUT[v]
}

// To8Bit converts a linear value to an 8-bit Adobe RGB encoded value, clipping
// the linear value to between 0.0 and 1.0.
//
// This implementation uses a fast look-up table and is approximate. For more
// accuracy, see ConvertLinearTo8Bit.
func To8Bit(v float32) uint8 {
	return linearToEncoded8LUT[linear.NormalisedTo9Bit(v)]
}

// To16Bit converts a linear value to a 16-bit Adobe RGB encoded value, clipping
// the linear value to between 0.0 and 1.0.
//
// This implementation uses a fast look-up table and is approximate. For more
// accuracy, see ConvertLinearTo16Bit.
func To16Bit(v float32) uint16 {
	if linearToEncoded16LUT != nil {
		return linearToEncoded16LUT[linear.NormalisedTo16Bit(v)]
	}
	return to16BitAndInitLUT(v)
}

func to16BitAndInitLUT(v float32) uint16 {
	initTo16BitLUTOnce.Do(func() {
		to16BitLUT := lut.BuildLinearTo16Bit(linearToEncoded)
		linearToEncoded16LUT = to16BitLUT[:]
	})
	return linearToEncoded16LUT[linear.NormalisedTo16Bit(v)]
}
