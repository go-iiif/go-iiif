package ciexyz

import (
	"github.com/mandykoh/prism/ciexyy"
	"github.com/mandykoh/prism/matrix"
)

var bradfordForward = matrix.Matrix3{
	{0.8951000, -0.7502000, 0.0389000},
	{0.2664000, 1.7135000, -0.0685000},
	{-0.1614000, 0.0367000, 1.0296000},
}

var bradfordInverse = bradfordForward.Inverse()

// ChromaticAdaptation represents an adaptation from one reference white point
// to another in XYZ space.
type ChromaticAdaptation matrix.Matrix3

// Apply transforms the specified colour using this chromatic adaptation.
func (ca ChromaticAdaptation) Apply(c Color) Color {
	return ColorFromV(matrix.Matrix3(ca).MulV(c.ToV()))
}

// AdaptBetweenXYYWhitePoints returns a ChromaticAdaptation from the source
// white point to the destination.
func AdaptBetweenXYYWhitePoints(srcWhite ciexyy.Color, dstWhite ciexyy.Color) ChromaticAdaptation {
	return AdaptBetweenXYZWhitePoints(ColorFromXYY(srcWhite), ColorFromXYY(dstWhite))
}

// AdaptBetweenXYZWhitePoints returns a ChromaticAdaptation from the source
// white point to the destination.
func AdaptBetweenXYZWhitePoints(srcWhite Color, dstWhite Color) ChromaticAdaptation {
	srcV := srcWhite.ToV()
	srcCSP := bradfordForward.MulV(srcV)

	dstV := dstWhite.ToV()
	dstCSP := bradfordForward.MulV(dstV)

	m := matrix.Matrix3{
		{dstCSP[0] / srcCSP[0], 0, 0},
		{0, dstCSP[1] / srcCSP[1], 0},
		{0, 0, dstCSP[2] / srcCSP[2]},
	}

	return ChromaticAdaptation(bradfordInverse.MulM(m).MulM(bradfordForward))
}
