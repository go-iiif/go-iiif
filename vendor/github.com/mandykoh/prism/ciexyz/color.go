package ciexyz

import (
	"github.com/mandykoh/prism/cielab"
	"github.com/mandykoh/prism/ciexyy"
	"github.com/mandykoh/prism/matrix"
	"math"
)

const constantE = 216.0 / 24389.0
const constantK = 24389.0 / 27.0

// Color represents a linear normalised colour in CIE XYZ space.
type Color struct {
	X float32
	Y float32
	Z float32
}

// ToLAB converts this colour to a CIE Lab colour given a reference white point.
func (c Color) ToLAB(whitePoint Color) cielab.Color {
	fx := componentToLAB(c.X, whitePoint.X)
	fy := componentToLAB(c.Y, whitePoint.Y)
	fz := componentToLAB(c.Z, whitePoint.Z)

	return cielab.Color{
		L: float32(116*fy - 16),
		A: float32(500 * (fx - fy)),
		B: float32(200 * (fy - fz)),
	}
}

// ToV returns this CIE XYZ colour as a vector.
func (c Color) ToV() matrix.Vector3 {
	return matrix.Vector3{float64(c.X), float64(c.Y), float64(c.Z)}
}

// ColorFromLAB creates a CIE XYZ Color instance from a CIE LAB representation
// given a reference white point.
func ColorFromLAB(lab cielab.Color, whitePoint Color) Color {
	fy := (float64(lab.L) + 16) / 116
	fx := float64(lab.A)/500 + fy
	fz := fy - float64(lab.B)/200

	xr := componentFromLAB(fx)
	zr := componentFromLAB(fz)

	var yr float64
	if lab.L > constantK*constantE {
		yr = math.Pow((float64(lab.L)+16)/116, 3)
	} else {
		yr = float64(lab.L) / constantK
	}

	return Color{
		X: float32(xr * float64(whitePoint.X)),
		Y: float32(yr * float64(whitePoint.Y)),
		Z: float32(zr * float64(whitePoint.Z)),
	}
}

// ColorFromV creates a CIE XYZ colour instance from a vector.
func ColorFromV(v matrix.Vector3) Color {
	return Color{
		X: float32(v[0]),
		Y: float32(v[1]),
		Z: float32(v[2]),
	}
}

// ColorFromXYY creates a CIE XYZ Color instance from a CIE xyY representation.
func ColorFromXYY(c ciexyy.Color) Color {
	return Color{
		X: c.X * c.YY / c.Y,
		Y: c.YY,
		Z: (1 - c.X - c.Y) * c.YY / c.Y,
	}
}
