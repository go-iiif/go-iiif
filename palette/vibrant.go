package palette

// SEE NOTES IN palette.go (20180601/thisisaaronland)

import (
	"fmt"
	"github.com/RobCherry/vibrant"
	"github.com/pwaller/go-hexcolor"
	"golang.org/x/image/draw"
	"image"
	"sort"
)

type VibrantPalette struct {
	Palette
	max_colors uint32
}

func NewVibrantPalette() (Palette, error) {

	v := VibrantPalette{
		max_colors: 24,
	}

	return &v, nil
}

func (v *VibrantPalette) Extract(im image.Image, limit int) ([]Color, error) {

	pb := vibrant.NewPaletteBuilder(im)
	// pb = pb.MaximumColorCount(v.max_colors)
	pb = pb.Scaler(draw.ApproxBiLinear)

	palette := pb.Generate()

	swatches := palette.Swatches()
	sort.Sort(populationSwatchSorter(swatches))

	colours := make([]Color, 0)

	for _, sw := range swatches {

		rgba := sw.RGBAInt()
		r, g, b, a := rgba.RGBA()

		hex := hexcolor.RGBAToHex(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))

		c := Color{
			Color: fmt.Sprintf("%s", hex),
		}

		colours = append(colours, c)

		if limit > 0 && len(colours) == limit {
			break
		}
	}

	return colours, nil
}

type populationSwatchSorter []*vibrant.Swatch

func (p populationSwatchSorter) Len() int           { return len(p) }
func (p populationSwatchSorter) Less(i, j int) bool { return p[i].Population() > p[j].Population() }
func (p populationSwatchSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
