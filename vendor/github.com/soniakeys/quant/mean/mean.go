// Copyright 2013 Sonia Keys.
// Licensed under MIT license.  See "license" file in this source tree.

// Mean is a simple color quantizer.  The algorithm successively divides the
// color space much like a median cut algorithm, but a mean statistic is used
// rather than a median.  In another simplification, there is no priority
// queue to order color blocks; linear search is used instead.
//
// An added sopphistication though, is that division proceeds in two stages,
// with somewhat different criteria used for the earlier cuts than for the
// later cuts.
//
// Motivation for using the mean is the observation that in a two stage
// algorithm, cuts are offset from the computed average so having the logically
// "correct" value of the median must not be that important.  Motivation
// for the linear search is that the number of blocks to search is limited
// to the target number of colors in the palette, which is small and typically
// limited to 256.  If n is 256, O(log n) and O(n) both become O(1).
package mean

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/soniakeys/quant"
	"github.com/soniakeys/quant/internal"
)

// Quantizer methods implement mean cut color quantization.
//
// The value is the target number of colors.
// Methods do not require pointer receivers, simply construct Quantizer
// objects with a type conversion.
//
// The type satisfies both quant.Quantizer and draw.Quantizer interfaces.
type Quantizer int

var _ quant.Quantizer = Quantizer(0)
var _ draw.Quantizer = Quantizer(0)

// Paletted performs color quantization and returns a paletted image.
//
// Returned is a new image.Paletted with no more than q colors.  Note though
// that image.Paletted is limited to 256 colors.
func (q Quantizer) Paletted(img image.Image) *image.Paletted {
	n := int(q)
	if n > 256 {
		n = 256
	}
	qz := newQuantizer(img, n)
	if n > 1 {
		qz.cluster() // cluster pixels by color
	}
	return qz.paletted() // generate paletted image from clusters
}

// Palette performs color quantization and returns a quant.Palette object.
//
// Returned is a palette with no more than q colors.  Q may be > 256.
func (q Quantizer) Palette(img image.Image) quant.Palette {
	qz := newQuantizer(img, int(q))
	if q > 1 {
		qz.cluster() // cluster pixels by color
	}
	return qz.palette()
}

// Quantize performs color quantization and returns a color.Palette.
//
// Following the behavior documented with the draw.Quantizer interface,
// "Quantize appends up to cap(p) - len(p) colors to p and returns the
// updated palette...."  This method does not limit the number of colors
// to 256.  Cap(p) or the quantity cap(p) - len(p) may be > 256.
// Also for this method the value of the Quantizer object is ignored.
func (Quantizer) Quantize(p color.Palette, m image.Image) color.Palette {
	n := cap(p) - len(p)
	qz := newQuantizer(m, n)
	if n > 1 {
		qz.cluster() // cluster pixels by color
	}
	return p[:len(p)+copy(p[len(p):cap(p)], qz.palette().ColorPalette())]
}

type quantizer struct {
	img image.Image // original image
	cs  []cluster   // len(cs) is the desired number of colors

	pxRGBA func(x, y int) (r, g, b, a uint32) // function to get original image RGBA color values
}

type point struct{ x, y int32 }

type cluster struct {
	px []point // list of points in the cluster
	// rgb const identifying dimension in color space with widest range
	widestDim int
	min, max  uint32 // min, max color values in dimension with widest range
	volume    uint64 // color volume
	priority  int    // early: population, late: population*volume
}

// indentifiers for RGB channels, or dimensions or axes of RGB color space
const (
	rgbR = iota
	rgbG
	rgbB
)

func newQuantizer(img image.Image, n int) *quantizer {
	if n < 1 {
		return &quantizer{img: img, pxRGBA: internal.PxRGBAfunc(img)}
	}
	// Make list of all pixels in image.
	b := img.Bounds()
	px := make([]point, (b.Max.X-b.Min.X)*(b.Max.Y-b.Min.Y))
	i := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			px[i].x = int32(x)
			px[i].y = int32(y)
			i++
		}
	}
	// Make clusters, populate first cluster with complete pixel list.
	cs := make([]cluster, n)
	cs[0].px = px
	return &quantizer{img: img, cs: cs, pxRGBA: internal.PxRGBAfunc(img)}
}

// Cluster by repeatedly splitting clusters in two stages.  For the first
// stage, prioritize by population and split tails off distribution in color
// dimension with widest range.  For the second stage, prioritize by the
// product of population and color volume, and split at the mean of the color
// values in the dimension with widest range.  Terminate when the desired number
// of clusters has been populated or when clusters cannot be further split.
func (qz *quantizer) cluster() {
	cs := qz.cs
	half := len(cs) / 2
	// cx is index of new cluster, populated at start of loop here, but
	// not yet analyzed.
	cx := 0
	c := &cs[cx]
	for {
		qz.setPriority(c, cx < half) // compute statistics for new cluster
		// determine cluster to split, sx
		sx := -1
		var maxP int
		for x := 0; x <= cx; x++ {
			// rule is to consider only clusters with non-zero color volume
			// and then split cluster with highest priority.
			if c := &cs[x]; c.max > c.min && c.priority > maxP {
				maxP = c.priority
				sx = x
			}
		}
		// If no clusters have any color variation, mark the end of the
		// cluster list and quit early.
		if sx < 0 {
			qz.cs = qz.cs[:cx+1]
			break
		}
		s := &cs[sx]
		m := qz.cutValue(s, cx < half) // get where to split cluster
		// point to next cluster to populate
		cx++
		c = &cs[cx]
		// populate c by splitting s into c and s at value m
		qz.split(s, c, m)
		// Normal exit is when all clusters are populated.
		if cx == len(cs)-1 {
			break
		}
		if cx == half {
			// change priorities on existing clusters
			for x := 0; x < cx; x++ {
				cs[x].priority =
					int(uint64(cs[x].priority) * (cs[x].volume >> 16) >> 29)
			}
		}
		qz.setPriority(s, cx < half) // set priority for newly split s
	}
}

func (q *quantizer) setPriority(c *cluster, early bool) {
	// Find extents of color values in each dimension.
	var maxR, maxG, maxB uint32
	minR := uint32(math.MaxUint32)
	minG := uint32(math.MaxUint32)
	minB := uint32(math.MaxUint32)
	for _, p := range c.px {
		r, g, b, _ := q.pxRGBA(int(p.x), int(p.y))
		if r < minR {
			minR = r
		}
		if r > maxR {
			maxR = r
		}
		if g < minG {
			minG = g
		}
		if g > maxG {
			maxG = g
		}
		if b < minB {
			minB = b
		}
		if b > maxB {
			maxB = b
		}
	}
	// See which color dimension had the widest range.
	w := rgbG
	min := minG
	max := maxG
	if maxR-minR > max-min {
		w = rgbR
		min = minR
		max = maxR
	}
	if maxB-minB > max-min {
		w = rgbB
		min = minB
		max = maxB
	}
	// store statistics
	c.widestDim = w
	c.min = min
	c.max = max
	c.volume = uint64(maxR-minR) * uint64(maxG-minG) * uint64(maxB-minB)
	c.priority = len(c.px)
	if !early {
		c.priority = int(uint64(c.priority) * (c.volume >> 16) >> 29)
	}
}

func (q *quantizer) cutValue(c *cluster, early bool) uint32 {
	var sum uint64
	switch c.widestDim {
	case rgbR:
		for _, p := range c.px {
			r, _, _, _ := q.pxRGBA(int(p.x), int(p.y))
			sum += uint64(r)
		}
	case rgbG:
		for _, p := range c.px {
			_, g, _, _ := q.pxRGBA(int(p.x), int(p.y))
			sum += uint64(g)
		}
	case rgbB:
		for _, p := range c.px {
			_, _, b, _ := q.pxRGBA(int(p.x), int(p.y))
			sum += uint64(b)
		}
	}
	mean := uint32(sum / uint64(len(c.px)))
	if early {
		// split in middle of longer tail rather than at mean
		if c.max-mean > mean-c.min {
			mean = (mean + c.max) / 2
		} else {
			mean = (mean + c.min) / 2
		}
	}
	return mean
}

func (q *quantizer) split(s, c *cluster, m uint32) {
	px := s.px
	var v uint32
	i := 0
	last := len(px) - 1
	for i <= last {
		// Get color value in appropriate dimension.
		r, g, b, _ := q.pxRGBA(int(px[i].x), int(px[i].y))
		switch s.widestDim {
		case rgbR:
			v = r
		case rgbG:
			v = g
		case rgbB:
			v = b
		}
		// Split into two non-empty parts at m.
		if v < m || m == s.min && v == m {
			i++
		} else {
			px[last], px[i] = px[i], px[last]
			last--
		}
	}
	// Split the pixel list.
	s.px = px[:i]
	c.px = px[i:]
}

func (qz *quantizer) paletted() *image.Paletted {
	cp := make(color.Palette, len(qz.cs))
	pi := image.NewPaletted(qz.img.Bounds(), cp)
	for i := range qz.cs {
		px := qz.cs[i].px
		// Average values in cluster to get palette color.
		var rsum, gsum, bsum int64
		for _, p := range px {
			r, g, b, _ := qz.pxRGBA(int(p.x), int(p.y))
			rsum += int64(r)
			gsum += int64(g)
			bsum += int64(b)
		}
		n64 := int64(len(px) << 8)
		cp[i] = color.RGBA{
			uint8(rsum / n64),
			uint8(gsum / n64),
			uint8(bsum / n64),
			0xff,
		}
		// set image pixels
		for _, p := range px {
			pi.SetColorIndex(int(p.x), int(p.y), uint8(i))
		}
	}
	return pi
}

func (qz *quantizer) palette() quant.Palette {
	cp := make(color.Palette, len(qz.cs))
	for i := range qz.cs {
		px := qz.cs[i].px
		// Average values in cluster to get palette color.
		var rsum, gsum, bsum int64
		for _, p := range px {
			r, g, b, _ := qz.pxRGBA(int(p.x), int(p.y))
			rsum += int64(r)
			gsum += int64(g)
			bsum += int64(b)
		}
		n64 := int64(len(px) << 8)
		cp[i] = color.RGBA{
			uint8(rsum / n64),
			uint8(gsum / n64),
			uint8(bsum / n64),
			0xff,
		}
	}
	return quant.LinearPalette{cp}
}
