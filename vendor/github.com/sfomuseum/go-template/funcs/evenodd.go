package funcs

import (
	"math"
)

func Mod(x float64, y float64) float64 {
	return math.Mod(x, y)
}

func IsOdd(i int) bool {

	if Mod(float64(i), 2.0) > 0.0 {
		return true
	}

	return false
}

func IsEven(i int) bool {

	if Mod(float64(i), 2.0) == 0.0 {
		return true
	}

	return false
}
