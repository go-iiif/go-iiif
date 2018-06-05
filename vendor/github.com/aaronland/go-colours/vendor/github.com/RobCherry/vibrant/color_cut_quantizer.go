package vibrant

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"
)

// Constants used by the Quantizer.
const (
	histogramSize = 1 << (quantizeWordWidth * 3)
)

// colorComponent represents a component of a color.
type colorComponent uint8

// Color components.
const (
	red colorComponent = iota
	green
	blue
)

// ColorCutQuantizer is a color quantizer based on the Median-cut algorithm, but optimized for picking out distinct
// colors rather than representation colors.
//
// The color space is represented as a 3-dimensional cube with each dimension being an RGB
// component. The cube is then repeatedly divided until we have reduced the color space to the
// requested number of colors. An average color is then generated from each cube.
//
// What makes this different to median-cut is that median-cut divided cubes so that all of the cubes
// have roughly the same population, where this quantizer divides boxes based on their color volume.
// This means that the color space is divided into distinct colors, rather than representative
// colors.
type ColorCutQuantizer struct {
	filters []Filter
}

// NewColorCutQuantizer creates a default ColorCutQuantizer.
func NewColorCutQuantizer() *ColorCutQuantizer {
	return &ColorCutQuantizer{
		filters: []Filter{DefaultFilter},
	}
}

// NewColorCutQuantizerWithFilters creates a ColorCutQuantizer with custom filters.
func NewColorCutQuantizerWithFilters(filters []Filter) *ColorCutQuantizer {
	return &ColorCutQuantizer{
		filters: filters,
	}
}

// Swatches returns a slice of swaches generated from the image.
func (q *ColorCutQuantizer) Swatches(colorCount uint32, m image.Image) []*Swatch {
	swatches := make([]*Swatch, 0, colorCount)

	if colorCount == 0 {
		return swatches
	}

	histogram := make([]uint32, histogramSize)

	imageBounds := m.Bounds()
	for y := imageBounds.Min.Y; y < imageBounds.Max.Y; y++ {
		for x := imageBounds.Min.X; x < imageBounds.Max.X; x++ {
			quantizedColor := QuantizedColorModel.Convert(m.At(x, y)).(QuantizedColor)
			histogram[quantizedColor]++
		}
	}

	// Now let's count the number of distinct colors
	distinctColorCount := uint32(0)
	for quantizedColorAsInt, count := range histogram {
		if count > 0 {
			quantizedColor := QuantizedColor(quantizedColorAsInt)
			if q.shouldIgnoreColor(quantizedColor) {
				// If we should ignore the color, set the population to 0
				histogram[quantizedColorAsInt] = 0
			} else {
				// If the color has population, increase the distinct color count
				distinctColorCount++
			}
		}
	}

	// Now lets go through create an array consisting of only distinct colors
	distinctQuantizedColors := make([]QuantizedColor, distinctColorCount)
	distinctColorIndex := 0
	for quantizedColorAsInt, count := range histogram {
		if count > 0 {
			distinctQuantizedColors[distinctColorIndex] = QuantizedColor(quantizedColorAsInt)
			distinctColorIndex++
		}
	}
	if distinctColorCount <= colorCount {
		for _, color := range distinctQuantizedColors {
			swatches = append(swatches, NewSwatch(color, histogram[color]))
		}
	} else {
		// We need to use quantization to reduce the number of colors.
		// Create a priority queue which is sorted by descending priority.
		// We will put in VBoxes prioritized by their volume.
		// This means we will always split the largest box in the queue.
		pq := newVBoxPriorityQueue(colorCount)
		// To start, offer a box which contains all of the colors
		pq.Offer(newVBox(distinctQuantizedColors, histogram, 0, uint32(len(distinctQuantizedColors))-1))
		// Now go through the boxes, splitting them until we have reached colorCount or there are no more boxes to split
		for uint32(pq.Len()) < colorCount {
			vbox := pq.Poll()
			if vbox.CanSplit() {
				// First split the box, and offer the result
				splitBox, _ := vbox.Split()
				pq.Offer(splitBox, vbox)
			} else {
				break
			}
		}
		// Finally, return the average colors of the color boxes
		for pq.Len() > 0 {
			swatch := pq.Poll().Swatch()
			// We're averaging a color box, so we can still get colors which we do not want, so we check again here.
			if !q.shouldIgnoreColor(swatch.Color()) {
				swatches = append(swatches, swatch)
			}
		}
	}
	return swatches
}

// Quantize populates a color palette for the image.
func (q *ColorCutQuantizer) Quantize(p color.Palette, m image.Image) color.Palette {
	numColors := cap(p) - len(p)
	if numColors <= 0 {
		return p
	}
	swatches := q.Swatches(uint32(numColors), m)
	for _, swatch := range swatches {
		p = append(p, swatch.Color())
	}
	return p
}

func (q *ColorCutQuantizer) shouldIgnoreColor(color color.Color) bool {
	for _, filter := range q.filters {
		if !filter.isAllowed(color) {
			return true
		}
	}
	return false
}

// vBox represents a tightly fitting box around a color space.
type vBox struct {
	colors    []QuantizedColor
	histogram []uint32

	lowerIndex, upperIndex                uint32
	minimumRed, minimumGreen, minimumBlue uint32
	maximumRed, maximumGreen, maximumBlue uint32
	population                            uint32
}

// newVBox creates a new vBox initialized with the provided QuantizedColor array, histogram, and indexes.
func newVBox(colors []QuantizedColor, histogram []uint32, lowerIndex uint32, upperIndex uint32) *vBox {
	vbox := &vBox{
		colors:     colors,
		histogram:  histogram,
		lowerIndex: lowerIndex,
		upperIndex: upperIndex,
	}
	vbox.fit()
	return vbox
}

// Volume returns the volume of the vBox.
func (v *vBox) Volume() uint32 {
	return (v.maximumRed - v.minimumRed + 1) * (v.maximumGreen - v.minimumGreen + 1) * (v.maximumBlue - v.minimumBlue + 1)
}

// CanSplit determines whether or not a vBox can be split.
func (v *vBox) CanSplit() bool {
	return v.Volume() > 1
}

// Split this color box at the mid-point along it's longest dimension
func (v *vBox) Split() (*vBox, error) {
	if !v.CanSplit() {
		return nil, fmt.Errorf("Can not split a box with only 1 color.")
	}
	splitIndex := v.findSplitIndex()
	splitVBox := newVBox(v.colors, v.histogram, splitIndex+1, v.upperIndex)
	v.upperIndex = splitIndex
	v.fit()
	return splitVBox, nil
}

// Fit the boundaries of this box to tightly fit the colors within the box.
func (v *vBox) fit() {
	var (
		localMinimumRed   uint8 = math.MaxUint8
		localMinimumGreen uint8 = math.MaxUint8
		localMinimumBlue  uint8 = math.MaxUint8
		localMaximumRed   uint8
		localMaximumGreen uint8
		localMaximumBlue  uint8
		count             uint32
	)
	for i := v.lowerIndex; i <= v.upperIndex; i++ {
		color := v.colors[i]
		count += v.histogram[color]
		r := color.QuantizedRed()
		g := color.QuantizedGreen()
		b := color.QuantizedBlue()
		if r < localMinimumRed {
			localMinimumRed = r
		}
		if r > localMaximumRed {
			localMaximumRed = r
		}
		if g < localMinimumGreen {
			localMinimumGreen = g
		}
		if g > localMaximumGreen {
			localMaximumGreen = g
		}
		if b < localMinimumBlue {
			localMinimumBlue = b
		}
		if b > localMaximumBlue {
			localMaximumBlue = b
		}
	}
	v.minimumRed = uint32(localMinimumRed)
	v.minimumGreen = uint32(localMinimumGreen)
	v.minimumBlue = uint32(localMinimumBlue)
	v.maximumRed = uint32(localMaximumRed)
	v.maximumGreen = uint32(localMaximumGreen)
	v.maximumBlue = uint32(localMaximumBlue)
	v.population = count
}

func (v *vBox) longestColorComponent() colorComponent {
	redLength := v.maximumRed - v.minimumRed
	greenLength := v.maximumGreen - v.minimumGreen
	blueLength := v.maximumBlue - v.minimumBlue
	if redLength >= greenLength && redLength >= blueLength {
		return red
	} else if greenLength >= redLength && greenLength >= blueLength {
		return green
	} else {
		return blue
	}
}

// Finds the point within this box's lowerIndex and upperIndex index of where to split.
//
// This is calculated by finding the longest color dimension, and then sorting the
// sub-array based on that dimension value in each color. The colors are then iterated over
// until a color is found with at least the midpoint of the whole box's dimension midpoint.
func (v *vBox) findSplitIndex() uint32 {
	longestColorComponent := v.longestColorComponent()

	colorSlice := v.colors[v.lowerIndex : v.upperIndex+1]
	// We need to sort the colors in this box based on the longest color dimension.
	if longestColorComponent != red {
		for i, color := range colorSlice {
			if longestColorComponent == green {
				colorSlice[i] = color.SwapRedGreen()
			} else {
				colorSlice[i] = color.SwapRedBlue()
			}
		}
	}
	sort.Sort(QuantizedColorSlice(colorSlice))
	// Now revert all of the colors so that they are packed as RGB again
	if longestColorComponent != red {
		for i, color := range colorSlice {
			if longestColorComponent == green {
				colorSlice[i] = color.SwapRedGreen()
			} else {
				colorSlice[i] = color.SwapRedBlue()
			}
		}
	}

	// Do not split a color across VBoxes.  This finds the highest index to search, leaving at least one color in the box.
	upperIndex := v.upperIndex - 1
	for upperIndex > v.lowerIndex {
		if v.colors[upperIndex] != v.colors[v.upperIndex] {
			break
		}
		upperIndex--
	}

	i := v.lowerIndex
	count := uint32(0)
	midPoint := v.population / 2
	for ; i < upperIndex; i++ {
		color := v.colors[i]
		count += v.histogram[color]
		if count >= midPoint {
			for ; color == v.colors[i+1]; i++ {
				// Continue on to the next color so that we do not split a color across VBoxes...
			}
			break
		}
	}
	return i
}

// Swatch generates a Swatch for the average color of this box.
func (v *vBox) Swatch() *Swatch {
	var (
		totalRed   uint32
		totalGreen uint32
		totalBlue  uint32
	)
	for _, color := range v.colors[v.lowerIndex : v.upperIndex+1] {
		colorPopulation := v.histogram[color]
		totalRed += colorPopulation * uint32(color.QuantizedRed())
		totalGreen += colorPopulation * uint32(color.QuantizedGreen())
		totalBlue += colorPopulation * uint32(color.QuantizedBlue())
	}
	totalRed = totalRed << (8 - quantizeWordWidth)
	totalGreen = totalGreen << (8 - quantizeWordWidth)
	totalBlue = totalBlue << (8 - quantizeWordWidth)
	averageRed := uint8(roundFloat64(float64(totalRed) / float64(v.population)))
	averageGreen := uint8(roundFloat64(float64(totalGreen) / float64(v.population)))
	averageBlue := uint8(roundFloat64(float64(totalBlue) / float64(v.population)))
	return NewSwatch(color.NRGBA{averageRed, averageGreen, averageBlue, 0xFF}, v.population)
}

type colorCutVBoxPriorityQueue struct {
	priorityQueue PriorityQueue
}

func newVBoxPriorityQueue(capacity uint32) colorCutVBoxPriorityQueue {
	return colorCutVBoxPriorityQueue{
		NewPriorityQueue(capacity, func(vbox interface{}) uint32 {
			return vbox.(*vBox).Volume()
		}),
	}
}

func (q colorCutVBoxPriorityQueue) Offer(items ...*vBox) {
	for _, item := range items {
		q.priorityQueue.Offer(item)
	}
}

func (q colorCutVBoxPriorityQueue) Poll() *vBox {
	return q.priorityQueue.Poll().(*vBox)
}

func (q colorCutVBoxPriorityQueue) Len() int {
	return q.priorityQueue.Len()
}
