package ciexyz

import (
	"github.com/mandykoh/prism/ciexyy"
	"github.com/mandykoh/prism/matrix"
	"math"
)

var D50 = Color{0.9642, 1.0, 0.8251}
var D65 = Color{0.95047, 1.0, 1.08883}

func componentFromLAB(f float64) float64 {
	if f3 := math.Pow(f, 3); f3 > constantE {
		return f3
	}
	return (116*f - 16) / constantK
}

func componentToLAB(v float32, wp float32) float64 {
	r := float64(v) / float64(wp)
	if r > constantE {
		return math.Pow(r, 1.0/3.0)
	}
	return (constantK*r + 16) / 116.0
}

// TransformFromXYZForXYYPrimaries generates the column matrix for converting
// colour values from CIE XYZ to a space defined by three RGB primary
// chromaticities and a reference white.
func TransformFromXYZForXYYPrimaries(r, g, b ciexyy.Color, whitePoint ciexyy.Color) matrix.Matrix3 {
	return TransformToXYZForXYYPrimaries(r, g, b, whitePoint).Inverse()
}

// TransformToXYZForXYYPrimaries generates the column matrix for converting
// colour values from a space defined by three primary RGB chromaticities and a
// reference white to CIE XYZ.
func TransformToXYZForXYYPrimaries(r, g, b ciexyy.Color, whitePoint ciexyy.Color) matrix.Matrix3 {
	m := matrix.Matrix3{
		ColorFromXYY(r).ToV(),
		ColorFromXYY(g).ToV(),
		ColorFromXYY(b).ToV(),
	}

	s := m.Inverse().MulV(ColorFromXYY(whitePoint).ToV())

	return matrix.Matrix3{
		m[0].MulS(s[0]),
		m[1].MulS(s[1]),
		m[2].MulS(s[2]),
	}
}
