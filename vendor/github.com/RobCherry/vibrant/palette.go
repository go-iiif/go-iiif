package vibrant

import (
	"golang.org/x/image/draw"
	"image"
	"log"
)

// Defaults used by the default PaletteBuilder.
const (
	DefaultResizeArea        = uint64(320 * 320)
	DefaultMaximumColorCount = uint32(64)
)

// Defaults used by the default PaletteBuilder.
var (
	DefaultScaler = draw.ApproxBiLinear
)

// Palette extracts prominent colors from an image.
//
// A number of colors with different profiles are extracted from the image:
//
//   Vibrant
//   Vibrant Dark
//   Vibrant Light
//   Muted
//   Muted Dark
//   Muted Light
//
// These can be retrieved via the appropriate getter method.
//
// Instances are created with a PaletteBuilder which supports several options to tweak the generated Palette.
type Palette struct {
	swatches         []*Swatch
	targets          []*Target
	selectedSwatches map[*Target]*Swatch
	usedColors       map[uint32]bool
	maxPopulation    uint32
}

func newPalette(swatches []*Swatch, targets []*Target) *Palette {
	maxPopulation := uint32(0)
	for _, swatch := range swatches {
		if swatch.Population() > maxPopulation {
			maxPopulation = swatch.Population()
		}
	}
	return &Palette{
		swatches:         swatches,
		targets:          targets,
		selectedSwatches: make(map[*Target]*Swatch),
		usedColors:       make(map[uint32]bool),
		maxPopulation:    maxPopulation,
	}
}

// Returns all of the swatches which make up the palette.
func (p *Palette) Swatches() []*Swatch {
	return p.swatches
}

// Returns the targets used to generate this palette.
func (p *Palette) Targets() []*Target {
	return p.targets
}

// Returns the most vibrant swatch in the palette. Might be nil.
func (p *Palette) VibrantSwatch() *Swatch {
	return p.SwatchForTarget(Vibrant)
}

// Returns a light and vibrant swatch from the palette. Might be nil.
func (p *Palette) LightVibrantSwatch() *Swatch {
	return p.SwatchForTarget(LightVibrant)
}

// Returns a dark and vibrant swatch from the palette. Might be nil.
func (p *Palette) DarkVibrantSwatch() *Swatch {
	return p.SwatchForTarget(DarkVibrant)
}

// Returns a muted swatch from the palette. Might be nil.
func (p *Palette) MutedSwatch() *Swatch {
	return p.SwatchForTarget(Muted)
}

// Returns a muted and light swatch from the palette. Might be nil.
func (p *Palette) LightMutedSwatch() *Swatch {
	return p.SwatchForTarget(LightMuted)
}

// Returns a muted and dark swatch from the palette. Might be nil.
func (p *Palette) DarkMutedSwatch() *Swatch {
	return p.SwatchForTarget(DarkMuted)
}

// Returns the most vibrant color in the palette as an RGB packed int.
func (p *Palette) VibrantColor(defaultColor uint32) uint32 {
	return p.ColorForTarget(Vibrant, defaultColor)
}

// Returns a light and vibrant color from the palette as an RGB packed int.
func (p *Palette) LightVibrantColor(defaultColor uint32) uint32 {
	return p.ColorForTarget(LightVibrant, defaultColor)
}

// Returns a dark and vibrant color from the palette as an RGB packed int.
func (p *Palette) DarkVibrantColor(defaultColor uint32) uint32 {
	return p.ColorForTarget(DarkVibrant, defaultColor)
}

// Returns a muted color from the palette as an RGB packed int.
func (p *Palette) MutedColor(defaultColor uint32) uint32 {
	return p.ColorForTarget(Muted, defaultColor)
}

// Returns a muted and light color from the palette as an RGB packed int.
func (p *Palette) LightMutedColor(defaultColor uint32) uint32 {
	return p.ColorForTarget(LightMuted, defaultColor)
}

// Returns a muted and dark color from the palette as an RGB packed int.
func (p *Palette) DarkMutedColor(defaultColor uint32) uint32 {
	return p.ColorForTarget(DarkMuted, defaultColor)
}

// Returns the selected swatch for the given target from the palette, or nil if one
// could not be found.
func (p *Palette) SwatchForTarget(target *Target) *Swatch {
	return p.selectedSwatches[target]
}

// Returns the selected color for the given target from the palette as an RGB packed int.
func (p *Palette) ColorForTarget(target *Target, defaultColor uint32) uint32 {
	swatch := p.SwatchForTarget(target)
	if swatch != nil {
		return swatch.RGBAInt().PackedRGB()
	}
	return defaultColor
}

func (p *Palette) generate() {
	// We need to make sure that the scored targets are generated first. This is so that
	// inherited targets have something to inherit from
	for _, target := range p.targets {
		p.selectedSwatches[target] = p.generateScoredTarget(target)
	}
	// We now clear out the used colors
	p.usedColors = make(map[uint32]bool)
}

func (p *Palette) generateScoredTarget(target *Target) *Swatch {
	maxScoreSwatch := p.maxScoredSwatchForTarget(target)
	if maxScoreSwatch != nil && target.IsExclusive() {
		// If we have a swatch, and the target is exclusive, add the color to the used list
		p.usedColors[maxScoreSwatch.RGBAInt().PackedRGB()] = true
	}
	return maxScoreSwatch
}

func (p *Palette) maxScoredSwatchForTarget(target *Target) *Swatch {
	scorer := target.Scorer(p.swatches)
	maxScore := 0.0
	var maxScoreSwatch *Swatch
	for _, swatch := range p.swatches {
		if p.shouldBeScoredForTarget(swatch, target) {
			score := scorer.Score(swatch)
			if maxScoreSwatch == nil || score > maxScore {
				maxScoreSwatch = swatch
				maxScore = score
			}
		}
	}
	return maxScoreSwatch
}

func (p *Palette) shouldBeScoredForTarget(swatch *Swatch, target *Target) bool {
	// Check whether the HSL values are within the correct ranges, and this color hasn't been used yet.
	hsl := swatch.HSL()
	return !p.usedColors[swatch.RGBAInt().PackedRGB()] &&
		hsl.S >= target.MinimumSaturation() && hsl.S <= target.MaximumSaturation() &&
		hsl.L >= target.MinimumLightness() && hsl.L <= target.MaximumLightness()
}

// PaletteBuilder is used for generating Palette instances.
type PaletteBuilder struct {
	image             image.Image
	region            image.Rectangle
	swatches          []*Swatch
	targets           []*Target
	filters           []Filter
	maximumColorCount uint32
	resizeArea        uint64
	scaler            draw.Scaler
}

// NewPaletteBuilder creates a new PaletteBuilder for an image.
func NewPaletteBuilder(image image.Image) *PaletteBuilder {
	return &PaletteBuilder{
		image:             image,
		region:            image.Bounds(),
		swatches:          make([]*Swatch, 0, DefaultMaximumColorCount),
		targets:           []*Target{Vibrant, Muted, LightVibrant, LightMuted, DarkVibrant, DarkMuted},
		filters:           []Filter{DefaultFilter},
		maximumColorCount: DefaultMaximumColorCount,
		resizeArea:        DefaultResizeArea,
		scaler:            DefaultScaler,
	}
}

// Set the maximum number of colors to use in the quantization step.
//
// Good values for depend on the source image type. For landscapes, good values are in
// the range 10-16. For images which are largely made up of people's faces then this
// value should be increased to ~24.
func (b *PaletteBuilder) MaximumColorCount(colors uint32) *PaletteBuilder {
	b.maximumColorCount = colors
	return b
}

// Set the resize value. If the image's area is greater than the value specified, then
// the image will be resized so that it's area matches the provided area. If the
// image is smaller or equal, the original is used as-is.
//
// This value has a large effect on the processing time. The larger the resized image is,
// the greater time it will take to generate the palette. The smaller the image is, the
// more detail is lost in the resulting image and thus less precision for color selection.
//
// A value of 0 can be used to disable resizing.
func (b *PaletteBuilder) ResizeImageArea(area uint64) *PaletteBuilder {
	b.resizeArea = area
	return b
}

// Specify the scaling function used to resize an image.  Set to nil to disable resizing.
func (b *PaletteBuilder) Scaler(scaler draw.Scaler) *PaletteBuilder {
	b.scaler = scaler
	return b
}

// Set a region of the image to be used exclusively when calculating the palette.
func (b *PaletteBuilder) Region(region image.Rectangle) *PaletteBuilder {
	if b.image != nil {
		b.region = region
		if !b.region.In(b.image.Bounds()) {
			log.Panicln("The given region must be within the image's dimensions.")
		}
	}
	return b
}

// Clear a previously specified region.
func (b *PaletteBuilder) ClearRegion() *PaletteBuilder {
	b.region = b.image.Bounds()
	return b
}

// Add a filter to be able to have fine grained control over which colors are
// allowed in the resulting palette.
func (b *PaletteBuilder) AddFilter(filter Filter) *PaletteBuilder {
	if filter != nil {
		b.filters = append(b.filters, filter)
	}
	return b
}

// Clear all added filters. This includes any default filters added automatically.
func (b *PaletteBuilder) ClearFilters() *PaletteBuilder {
	b.filters = make([]Filter, 0)
	return b
}

// Add a target profile to be generated in the palette.
//
// You can retrieve the result via Palette#getSwatchForTarget(Target).
func (b *PaletteBuilder) AddTarget(target *Target) *PaletteBuilder {
	shouldAppend := true
	for _, t := range b.targets {
		if t == target {
			shouldAppend = false
			break
		}
	}
	if shouldAppend {
		b.targets = append(b.targets, target)
	}
	return b
}

// Clear all added targets. This includes any default targets added automatically.
func (b *PaletteBuilder) ClearTargets() *PaletteBuilder {
	b.targets = make([]*Target, 0)
	return b
}

// Generate and return the Palette.
func (b *PaletteBuilder) Generate() *Palette {
	if b.image != nil {
		// Select the region of the image if possible...
		m := b.image
		if _, ok := b.image.(SubImager); ok {
			m = m.(SubImager).SubImage(b.region)
		}

		// Next we scale the image if necessary...
		m = ScaleImageDown(m, b.resizeArea, b.scaler)

		// Now generate swatches from the image...
		quantizer := NewColorCutQuantizerWithFilters(b.filters)
		b.swatches = quantizer.Swatches(b.maximumColorCount, m)
	}
	palette := newPalette(b.swatches, b.targets)
	palette.generate()
	return palette
}
