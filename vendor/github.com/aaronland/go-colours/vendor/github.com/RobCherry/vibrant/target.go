package vibrant

import (
	"fmt"
	"math"
)

// Constants used by the default targets.
const (
	MinNormalLightness    = 0.3
	TargetNormalLightness = 0.5
	MaxNormalLightness    = 0.7

	MinLightLightness    = 0.55
	TargetLightLightness = 0.74

	TargetDarkLightness = 0.26
	MaxDarkLightness    = 0.45

	TargetVibrantSaturation = 1
	MinVibrantSaturation    = 0.35

	TargetMutedSaturation = 0.3
	MaxMutedSaturation    = 0.4

	// Weights used by the default scoring function.
	DefaultSaturationWeight = 0.375
	DefaultLightnessWeight  = 0.575
	DefaultPopulationWeight = 0.05
)

// The default targets.
var (
	// A target which has the characteristics of a vibrant color which is light in luminance.
	LightVibrant *Target
	// A target which has the characteristics of a vibrant color which is neither light or dark.
	Vibrant *Target
	// A target which has the characteristics of a vibrant color which is dark in luminance.
	DarkVibrant *Target
	// A target which has the characteristics of a muted color which is light in luminance.
	LightMuted *Target
	// A target which has the characteristics of a muted color which is neither light or dark.
	Muted *Target
	// A target which has the characteristics of a muted color which is dark in luminance.
	DarkMuted *Target
)

// The default scorer factory.
var (
	// A function which generates a Scorer a target.
	DefaultScorerFactory ScorerFactory
)

func init() {
	DefaultScorerFactory = ScorerFactory(func(target Target, swatches []*Swatch) Scorer {
		maxPopulation := uint32(0)
		for _, swatch := range swatches {
			if swatch.Population() > maxPopulation {
				maxPopulation = swatch.Population()
			}
		}
		return NewDefaultScorer(target, maxPopulation)
	})

	LightVibrant = NewTargetBuilder().
		MinimumLightness(MinLightLightness).
		TargetLightness(TargetLightLightness).
		MinimumSaturation(MinVibrantSaturation).
		TargetSaturation(TargetVibrantSaturation).
		Build()

	Vibrant = NewTargetBuilder().
		MinimumLightness(MinNormalLightness).
		TargetLightness(TargetNormalLightness).
		MaximumLightness(MaxNormalLightness).
		MinimumSaturation(MinVibrantSaturation).
		TargetSaturation(TargetVibrantSaturation).
		Build()

	DarkVibrant = NewTargetBuilder().
		TargetLightness(TargetDarkLightness).
		MaximumLightness(MaxDarkLightness).
		MinimumSaturation(MinVibrantSaturation).
		TargetSaturation(TargetVibrantSaturation).
		Build()

	LightMuted = NewTargetBuilder().
		MinimumLightness(MinLightLightness).
		TargetLightness(TargetLightLightness).
		TargetSaturation(TargetMutedSaturation).
		MaximumSaturation(MaxMutedSaturation).
		Build()

	Muted = NewTargetBuilder().
		MinimumLightness(MinNormalLightness).
		TargetLightness(TargetNormalLightness).
		MaximumLightness(MaxNormalLightness).
		TargetSaturation(TargetMutedSaturation).
		MaximumSaturation(MaxMutedSaturation).
		Build()

	DarkMuted = NewTargetBuilder().
		TargetLightness(TargetDarkLightness).
		MaximumLightness(MaxDarkLightness).
		TargetSaturation(TargetMutedSaturation).
		MaximumSaturation(MaxMutedSaturation).
		Build()
}

// Scorer used to score a Swatch for a Target.
type Scorer interface {
	// Score returns the score for a Swatch.
	Score(*Swatch) float64
}

// DefaultScorer is the default scorer implementation used to score a Target.
type DefaultScorer struct {
	target        Target
	maxPopulation uint32
}

// NewDefaultScorer creates a new default scorer.
func NewDefaultScorer(target Target, maxPopulation uint32) Scorer {
	return &DefaultScorer{
		target:        target,
		maxPopulation: maxPopulation,
	}
}

// Score returns the score for a Swatch.
func (s DefaultScorer) Score(swatch *Swatch) float64 {
	hsl := swatch.HSL()
	if hsl.S < s.target.MinimumSaturation() || hsl.S > s.target.MaximumSaturation() {
		return 0
	}
	if hsl.L < s.target.MinimumLightness() || hsl.L > s.target.MaximumLightness() {
		return 0
	}
	saturationScore := DefaultSaturationWeight * (1.0 - math.Abs(hsl.S-s.target.TargetSaturation()))
	lightnessScore := DefaultLightnessWeight * math.Pow(1.0-math.Abs(hsl.L-s.target.TargetLightness()), 2)
	populationScore := DefaultPopulationWeight * (float64(swatch.Population()) / float64(s.maxPopulation))
	return saturationScore + lightnessScore + populationScore
}

// ScorerFactory creates a Scorer.
type ScorerFactory func(Target, []*Swatch) Scorer

// Target allows custom selection of colors in a Palette's generation. Instances can be created via TargetBuilder.
//
// To use the target, use the PaletteBuilder#addTarget(Target) API when building a Palette.
type Target struct {
	minimumSaturation float64
	targetSaturation  float64
	maximumSaturation float64

	minimumLightness float64
	targetLightness  float64
	maximumLightness float64

	isExclusive bool

	scorerFactory ScorerFactory
}

func newTarget() *Target {
	return &Target{
		minimumSaturation: 0,
		targetSaturation:  0.5,
		maximumSaturation: 1,
		minimumLightness:  0,
		targetLightness:   0.5,
		maximumLightness:  1,
		isExclusive:       true,
		scorerFactory:     DefaultScorerFactory,
	}
}

func newTargetFromTarget(other *Target) *Target {
	return &Target{
		minimumSaturation: other.minimumSaturation,
		targetSaturation:  other.targetSaturation,
		maximumSaturation: other.maximumSaturation,
		minimumLightness:  other.minimumLightness,
		targetLightness:   other.targetLightness,
		maximumLightness:  other.maximumLightness,
		isExclusive:       other.isExclusive,
		scorerFactory:     other.scorerFactory,
	}
}

func (t Target) String() string {
	return fmt.Sprintf("Saturation (%0.2f, %0.2f, %0.2f) Lightness (%0.2f, %0.2f, %0.2f) Exclusive (%v)", t.minimumSaturation, t.targetSaturation, t.maximumLightness, t.minimumLightness, t.targetLightness, t.maximumLightness, t.isExclusive)
}

// The minimum saturation value for this target.
func (t Target) MinimumSaturation() float64 {
	return t.minimumSaturation
}

// The target saturation value for this target.
func (t Target) TargetSaturation() float64 {
	return t.targetSaturation
}

// The maximum saturation value for this target.
func (t Target) MaximumSaturation() float64 {
	return t.maximumSaturation
}

// The minimum lightness value for this target.
func (t Target) MinimumLightness() float64 {
	return t.minimumLightness
}

// The target lightness value for this target.
func (t Target) TargetLightness() float64 {
	return t.targetLightness
}

// The maximum lightness value for this target.
func (t Target) MaximumLightness() float64 {
	return t.maximumLightness
}

// Returns whether any color selected for this target is exclusive for this target only.
//
// If false, then the color can be selected for other targets.
func (t Target) IsExclusive() bool {
	return t.isExclusive
}

// Returns a Scorer for the target.
func (t Target) Scorer(swatches []*Swatch) Scorer {
	return t.scorerFactory(t, swatches)
}

// TargetBuilder is used for generating Target instances.
type TargetBuilder struct {
	result *Target
}

// NewTargetBuilder creates a new TargetBuilder from scratch.
func NewTargetBuilder() *TargetBuilder {
	return &TargetBuilder{
		result: newTarget(),
	}
}

// NewTargetBuilderFromTarget creates a new TargetBuilder based on an existing Target.
func NewTargetBuilderFromTarget(target *Target) *TargetBuilder {
	return &TargetBuilder{
		result: newTargetFromTarget(target),
	}
}

// Set the minimum saturation value for this target.
func (b *TargetBuilder) MinimumSaturation(value float64) *TargetBuilder {
	b.result.minimumSaturation = clampFloat64(value, 0, 1)
	return b
}

// Set the target/ideal saturation value for this target.
func (b *TargetBuilder) TargetSaturation(value float64) *TargetBuilder {
	b.result.targetSaturation = clampFloat64(value, 0, 1)
	return b
}

// Set the maximum saturation value for this target.
func (b *TargetBuilder) MaximumSaturation(value float64) *TargetBuilder {
	b.result.maximumSaturation = clampFloat64(value, 0, 1)
	return b
}

// Set the minimum lightness value for this target.
func (b *TargetBuilder) MinimumLightness(value float64) *TargetBuilder {
	b.result.minimumLightness = clampFloat64(value, 0, 1)
	return b
}

// Set the target/ideal lightness value for this target.
func (b *TargetBuilder) TargetLightness(value float64) *TargetBuilder {
	b.result.targetLightness = clampFloat64(value, 0, 1)
	return b
}

// Set the maximum lightness value for this target.
func (b *TargetBuilder) MaximumLightness(value float64) *TargetBuilder {
	b.result.maximumLightness = clampFloat64(value, 0, 1)
	return b
}

// Set whether any color selected for this target is exclusive to this target only.
// Defaults to true.
func (b *TargetBuilder) IsExclusive(exclusive bool) *TargetBuilder {
	b.result.isExclusive = exclusive
	return b
}

// Set the ScorerFactory for the target.  The ScorerFactory is used to generate a Scorer for a target from a slice of swatches.
func (b *TargetBuilder) ScorerFactory(factory ScorerFactory) *TargetBuilder {
	b.result.scorerFactory = factory
	return b
}

// Build returns a copy of the Target.
func (b *TargetBuilder) Build() *Target {
	return newTargetFromTarget(b.result)
}
