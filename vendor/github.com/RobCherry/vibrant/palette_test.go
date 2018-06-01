package vibrant

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	_ "image/jpeg"
	"os"
	"sort"
	"strconv"
	"testing"
)

func ExamplePaletteBuilder() {
	file, _ := os.Open("test_files/1.jpg")
	decodedImage, _, _ := image.Decode(file)
	palette := NewPaletteBuilder(decodedImage).Generate()
	// Iterate over the swatches in the palette...
	for _, swatch := range palette.Swatches() {
		fmt.Printf("Swatch has color %v and population %d\n", swatch.RGBAInt(), swatch.Population())
	}
	for _, target := range palette.Targets() {
		_ = palette.SwatchForTarget(target)
		// Do something with the swatch for a given target...
	}
}

func ExamplePaletteBuilder_maximumColorCount() {
	file, _ := os.Open("test_files/1.jpg")
	decodedImage, _, _ := image.Decode(file)
	// Use a custom color count.
	palette := NewPaletteBuilder(decodedImage).MaximumColorCount(32).Generate()
	// Iterate over the swatches in the palette...
	for _, swatch := range palette.Swatches() {
		fmt.Printf("Swatch has color %v and population %d\n", swatch.RGBAInt(), swatch.Population())
	}
	for _, target := range palette.Targets() {
		_ = palette.SwatchForTarget(target)
		// Do something with the swatch for a given target...
	}
}

func ExamplePaletteBuilder_resizeImageArea() {
	file, _ := os.Open("test_files/1.jpg")
	decodedImage, _, _ := image.Decode(file)
	// Use a custom resize image area and scaler.
	palette := NewPaletteBuilder(decodedImage).ResizeImageArea(160 * 160).Scaler(draw.CatmullRom).Generate()
	// Iterate over the swatches in the palette...
	for _, swatch := range palette.Swatches() {
		fmt.Printf("Swatch has color %v and population %d\n", swatch.RGBAInt(), swatch.Population())
	}
	for _, target := range palette.Targets() {
		_ = palette.SwatchForTarget(target)
		// Do something with the swatch for a given target...
	}
}

func TestPalette_Swatches(t *testing.T) {
	tests := map[string]map[uint32]Uint32Slice{
		"test_files/1.jpg": {
			16: Uint32Slice([]uint32{
				0xFF2E222B,
				0xFF433534,
				0xFF6B433E,
				0xFF81635F,
				0xFF66505C,
				0xFF81859B,
				0xFF3B3951,
				0xFF7B5659,
				0xFFAB6A44,
				0xFFC7AF47,
				0xFF18BDE0,
				0xFFAD837F,
				0xFFC2C094,
				0xFF3A4D6E,
				0xFF246890,
				0xFF1F9A93,
			}),
			64: Uint32Slice([]uint32{
				0xFF2C1E1E,
				0xFF433534,
				0xFF2E2732,
				0xFF5A4444,
				0xFF66505C,
				0xFF7B5654,
				0xFF825E5B,
				0xFF69443C,
				0xFF816D8F,
				0xFF3F3740,
				0xFF383A62,
				0xFF322A38,
				0xFF966745,
				0xFFC16E44,
				0xFF886868,
				0xFF774436,
				0xFF3A4D6E,
				0xFF12E5EA,
				0xFFC4B873,
				0xFFAB7D70,
				0xFF957873,
				0xFF765866,
				0xFF252744,
				0xFFB7A136,
				0xFFB9A855,
				0xFF90392A,
				0xFF7786E1,
				0xFFD6CA4F,
				0xFFD7A944,
				0xFFAE887E,
				0xFF15B0E0,
				0xFF7F5862,
				0xFF766968,
				0xFF3A2261,
				0xFF432495,
				0xFF7A6D4D,
				0xFF826C5F,
				0xFF806868,
				0xFF77A581,
				0xFFB23633,
				0xFFD43B34,
				0xFF73F5E0,
				0xFF0F81CA,
				0xFF22686E,
				0xFF2668B2,
				0xFF1D858D,
				0xFF4376CF,
				0xFFAFB29C,
				0xFFCDBF9C,
				0xFFA9DCCD,
				0xFFDAD4CD,
				0xFFB28F8A,
				0xFF6BC3C5,
				0xFF72E3C9,
				0xFF1BB6AA,
				0xFF27A989,
				0xFFAB87AA,
				0xFFA780D8,
				0xFF805860,
				0xFF805868,
				0xFFAE3264,
				0xFFE94276,
				0xFFAF376D,
				0xFF9F49BB,
			}),
		},
		"test_files/2.jpg": {
			16: Uint32Slice([]uint32{
				0xFF20201D,
				0xFFDCDEE0,
				0xFF2C2D40,
				0xFF8F9095,
				0xFF0E146C,
				0xFFBEBFC0,
				0xFF4E525E,
				0xFF9B8452,
				0xFFAFAFAF,
				0xFFADAC9F,
				0xFFEEC60C,
				0xFF70695A,
				0xFF565E73,
				0xFF3A55B6,
				0xFF787665,
				0xFF955B1B,
			}),
			64: Uint32Slice([]uint32{
				0xFF0D0F1D,
				0xFFE8E8E6,
				0xFF11143F,
				0xFF494841,
				0xFF28211C,
				0xFF45471E,
				0xFFD8D8D7,
				0xFF828894,
				0xFF0A0E5F,
				0xFFB5B9BD,
				0xFF4E525E,
				0xFFC7C6C3,
				0xFFDBDEE0,
				0xFFAFAFAF,
				0xFFD2D0D0,
				0xFF847F6F,
				0xFF989594,
				0xFF0F156F,
				0xFFA6A8A4,
				0xFF787665,
				0xFFC5CFDB,
				0xFFF8D80C,
				0xFF6F6F72,
				0xFF726341,
				0xFF9E9B9A,
				0xFFC78907,
				0xFF686770,
				0xFF161C7B,
				0xFFC7C8D1,
				0xFFEDBE0A,
				0xFF354778,
				0xFF596573,
				0xFF93876F,
				0xFF2145D2,
				0xFF31409B,
				0xFFA79A94,
				0xFF506996,
				0xFFA9A09A,
				0xFF162089,
				0xFF1423A1,
				0xFFA7C0E9,
				0xFF995704,
				0xFFB3A899,
				0xFFE5A909,
				0xFF395DDD,
				0xFF5571CD,
				0xFF90A3E7,
				0xFF728CE1,
				0xFF9A9575,
				0xFF855E39,
				0xFFD0AC04,
				0xFFF1E77B,
				0xFFB9B494,
				0xFFAB8A36,
				0xFFCF9221,
				0xFF9F5F14,
				0xFFCAA521,
				0xFFD1BE23,
				0xFF8E6550,
				0xFFA9671D,
				0xFF996745,
				0xFFAB6428,
				0xFFC3506B,
				0xFFC85B38,
			}),
		},
		"test_files/3.jpg": {
			16: Uint32Slice([]uint32{
				0xFFF0EBD7,
				0xFF073125,
				0xFFE1D6BC,
				0xFF5F7A63,
				0xFF244B3A,
				0xFFD4C8AC,
				0xFF80927C,
				0xFFADAB8C,
				0xFF3E513F,
				0xFFCCBFA1,
				0xFF663723,
				0xFFC6B897,
				0xFF9E3E34,
				0xFFC1B094,
				0xFFC94E4E,
				0xFFD96863,
			}),
			64: Uint32Slice([]uint32{
				0xFFE1D6BC,
				0xFFF4F1E3,
				0xFF04251B,
				0xFF0A3D2F,
				0xFFF2E3CC,
				0xFFD0C9A9,
				0xFF657F68,
				0xFF4B6E59,
				0xFF245846,
				0xFFCCBFA1,
				0xFF547A67,
				0xFFE8E8C8,
				0xFFE8E8D5,
				0xFF6F8978,
				0xFFA4A081,
				0xFFE0C8B1,
				0xFFAEA88A,
				0xFF88947D,
				0xFFC6B897,
				0xFFB0B293,
				0xFF949B81,
				0xFF3C624F,
				0xFF758067,
				0xFFDFE0C8,
				0xFF799379,
				0xFFE7E1C6,
				0xFF1B4536,
				0xFFB8B898,
				0xFF254839,
				0xFF3A5645,
				0xFF612E19,
				0xFF2D2612,
				0xFF858364,
				0xFF723625,
				0xFFBBB499,
				0xFFC8AD90,
				0xFF243928,
				0xFF482816,
				0xFF423D28,
				0xFF873728,
				0xFFB74740,
				0xFFCA5556,
				0xFF908769,
				0xFFC64644,
				0xFF5C583F,
				0xFFA63836,
				0xFFD75A59,
				0xFF8A4732,
				0xFFA84038,
				0xFFA48567,
				0xFFCE4B4B,
				0xFFAA443B,
				0xFFE26567,
				0xFF9A906F,
				0xFF983830,
				0xFF9C3833,
				0xFFC49579,
				0xFF983828,
				0xFFE9706B,
				0xFFA12727,
				0xFF96372F,
				0xFFE97A78,
				0xFFDC9486,
				0xFFA4CAB4,
			}),
		},
		"test_files/4.jpg": {
			16: Uint32Slice([]uint32{
				0xFF151717,
				0xFF1B1923,
				0xFF181D2E,
				0xFF5B3F29,
				0xFFD0B15F,
				0xFF252937,
				0xFFB38B37,
				0xFFB79648,
				0xFF705538,
				0xFFB89454,
				0xFFCD6A89,
				0xFF875E42,
				0xFF765C53,
				0xFF9F6781,
				0xFF725968,
				0xFF647EAB,
			}),
			64: Uint32Slice([]uint32{
				0xFF151717,
				0xFF181825,
				0xFF181D2E,
				0xFF251816,
				0xFF351B18,
				0xFF451E19,
				0xFF402E29,
				0xFF562419,
				0xFF390A0C,
				0xFF7F5529,
				0xFF182038,
				0xFFC1A64F,
				0xFFB89538,
				0xFF94090F,
				0xFFD0B062,
				0xFFB79648,
				0xFF705538,
				0xFF7B5C1A,
				0xFF872519,
				0xFF67291B,
				0xFFA46A1B,
				0xFFB89454,
				0xFFDBBA6C,
				0xFF9D0C1B,
				0xFFAF883C,
				0xFF182838,
				0xFFD6A84E,
				0xFFD9BA48,
				0xFFB20E1F,
				0xFF8B5E31,
				0xFFE1C078,
				0xFFA88836,
				0xFFB16A30,
				0xFF97833F,
				0xFF773844,
				0xFFAE852A,
				0xFFE5C586,
				0xFFC89859,
				0xFF654157,
				0xFF8A7D4E,
				0xFFBD4273,
				0xFF9B9480,
				0xFFC09664,
				0xFFDF8792,
				0xFFEDC89E,
				0xFFEDC2D3,
				0xFFCED8F6,
				0xFFDB5ADB,
				0xFFCD2590,
				0xFFF8EBE8,
				0xFFB81894,
				0xFFA04E7A,
				0xFF8E5E66,
				0xFFA61A8C,
				0xFFD024A8,
				0xFF848E98,
				0xFFEB85D8,
				0xFFE33DE8,
				0xFFC82CB8,
				0xFFE874BC,
				0xFF404078,
				0xFF606894,
				0xFF2874E8,
				0xFF6C685C,
			}),
		},
	}
	for path, data := range tests {
		for maximumColorCount, expectedSwatches := range data {
			expectedSwatches.Sort()
			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			original, _, err := image.Decode(file)
			if err != nil {
				t.Fatal(err)
			}
			palette := NewPaletteBuilder(original).ResizeImageArea(0).MaximumColorCount(maximumColorCount).Generate()
			actualSwatches := palette.Swatches()
			if len(actualSwatches) != len(expectedSwatches) {
				t.Errorf("Expected %d swatches but generated %d for %s", len(expectedSwatches), len(actualSwatches), path)
			}
			for _, swatch := range actualSwatches {
				actual := swatch.RGBAInt().PackedRGBA()
				i := expectedSwatches.Search(actual)
				if !(i < len(expectedSwatches) && expectedSwatches[i] == actual) {
					t.Errorf("Swatch 0x%s was not expected in %s", strconv.FormatUint(uint64(actual), 16), path)
				}
			}
		}
	}
}

func TestPaletteBuilder_Region(t *testing.T) {
	file, err := os.Open("test_files/1.jpg")
	if err != nil {
		t.Fatal(err)
	}
	original, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	tests := map[image.Rectangle]RGBAInt{
		original.Bounds():              RGBAInt(0xFF05DDEC),
		image.Rect(205, 230, 260, 260): RGBAInt(0xFFF81454),
		image.Rect(340, 380, 375, 410): RGBAInt(0xFF08F8F8),
		image.Rect(570, 380, 600, 400): RGBAInt(0xFF483CC7),
	}
	for region, expected := range tests {
		actual := NewPaletteBuilder(original).Region(region).Generate().VibrantSwatch().RGBAInt()
		if actual != expected {
			t.Errorf("Expected %v but generated %v as the vibrant color in %v\n", expected, actual, region)
		}
	}
}

func SearchUint32s(a []uint32, x uint32) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

type Uint32Slice []uint32

func (p Uint32Slice) Len() int            { return len(p) }
func (p Uint32Slice) Less(i, j int) bool  { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p Uint32Slice) Sort()               { sort.Sort(p) }
func (p Uint32Slice) Search(x uint32) int { return SearchUint32s(p, x) }
