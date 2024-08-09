package matrix

// Vector3 represents a 3-element vector.
type Vector3 [3]float64

// MulS returns the result of multiplying this vector by a scalar.
func (v Vector3) MulS(s float64) Vector3 {
	return Vector3{v[0] * s, v[1] * s, v[2] * s}
}

// Dot returns the dot product between two vectors.
func Dot(v1, v2 Vector3) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}
