package vibrant

import (
	"math"
	"testing"
)

func TestDefaultScorer_Score(t *testing.T) {
	tests := map[*Swatch]float64{
		NewSwatch(HSL{0, 1, 0.29999, 1}, 100):   0,
		NewSwatch(HSL{0, 1, 0.3, 1}, 100):       0.793,
		NewSwatch(HSL{0, 1, 0.45, 1}, 100):      0.9439375,
		NewSwatch(HSL{0, 1, 0.55, 1}, 100):      0.9439375,
		NewSwatch(HSL{0, 1, 0.7, 1}, 100):       0.793,
		NewSwatch(HSL{0, 1, 0.70001, 1}, 100):   0,
		NewSwatch(HSL{0, 0.34999, 0.5, 1}, 100): 0,
		NewSwatch(HSL{0, 0.35, 0.5, 1}, 100):    0.75625,
		NewSwatch(HSL{0, 0.75, 0.5, 1}, 100):    0.90625,
		NewSwatch(HSL{0, 0.95, 0.5, 1}, 100):    0.98125,
		NewSwatch(HSL{0, 0.99, 0.5, 1}, 100):    0.99625,
		NewSwatch(HSL{0, 1, 0.5, 1}, 100):       1,
		NewSwatch(HSL{0, 1, 0.5, 1}, 50):        0.975,
		NewSwatch(HSL{0, 1, 0.5, 1}, 0):         0.95,
	}
	scorer := NewDefaultScorer(*Vibrant, 100)
	for swatch, expected := range tests {
		score := scorer.Score(swatch)
		if !inDelta(expected, score, 0.0000001) {
			t.Errorf("Expected score %v does not match %v\n", expected, score)
		}
	}
}

func inDelta(expected, actual, delta float64) bool {
	return math.IsNaN(expected) == math.IsNaN(actual) && math.Abs(expected-actual) < delta
}
