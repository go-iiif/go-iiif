package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/RobCherry/vibrant"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var scalerByName = map[string]draw.Scaler{
	"NearestNeighbor": draw.NearestNeighbor,
	"ApproxBiLinear":  draw.ApproxBiLinear,
	"BiLinear":        draw.BiLinear,
	"CatmullRom":      draw.CatmullRom,
}

var (
	inputFile         string
	outputFile        string
	maximumColorCount uint64 = 64
	resizeImageArea   uint64
	scalerName        string
	debug             bool
)

func init() {
	var b bytes.Buffer
	for name, s := range scalerByName {
		if s == vibrant.DefaultScaler {
			scalerName = name
		}
		if b.Len() > 0 {
			b.WriteString(", ")
		}
		b.WriteString(name)
	}
	availableScalers := b.String()

	flag.StringVar(&outputFile, "outputFile", "", "Output location for a quantized version of the image.")
	flag.StringVar(&outputFile, "o", outputFile, "Output location for a quantized version of the image.")

	flag.Uint64Var(&maximumColorCount, "maximumColorCount", maximumColorCount, "Maximum color count.")
	flag.Uint64Var(&maximumColorCount, "m", maximumColorCount, "Maximum color count.")

	flag.Uint64Var(&resizeImageArea, "resizeImageArea", resizeImageArea, "Resize image area.")
	flag.Uint64Var(&resizeImageArea, "r", resizeImageArea, "Resize image area.")

	flag.BoolVar(&debug, "debug", debug, "Debug mode.")
	flag.BoolVar(&debug, "d", debug, "Debug mode.")

	flag.StringVar(&scalerName, "scaler", scalerName, "Scaler.  One of: "+availableScalers)
	flag.StringVar(&scalerName, "s", scalerName, "Scaler.  One of: "+availableScalers)

	flag.Parse()

	inputFile = flag.Arg(0)
}

func main() {
	input, err := os.Open(inputFile)
	if os.IsNotExist(err) {
		fmt.Printf("Unable to open %s\n", inputFile)
		os.Exit(1)
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer input.Close()

	inputImage, _, err := image.Decode(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	paletteBuilder := vibrant.NewPaletteBuilder(inputImage).
		MaximumColorCount(uint32(maximumColorCount)).
		ResizeImageArea(resizeImageArea)

	validScaler := false
	for name, scaler := range scalerByName {
		if name == scalerName {
			validScaler = true
			paletteBuilder = paletteBuilder.Scaler(scaler)
			break
		}
	}
	if !validScaler {
		fmt.Printf("%s is not a valid scaler\n", scalerName)
		os.Exit(1)
	}

	start := time.Now()
	palette := paletteBuilder.Generate()
	elapsed := time.Since(start)
	if debug {
		fmt.Printf("Palette generation took approximately %s\n\n", elapsed)
	}

	swatches := palette.Swatches()
	sort.Sort(populationSwatchSorter(swatches))
	colorPalette := make(color.Palette, 0, len(swatches))
	for i, swatch := range swatches {
		colorPalette = append(colorPalette, swatch.Color())
		if debug {
			fmt.Printf("Swatch: %s (%d)\n", swatch.RGBAInt(), swatch.Population())
			fmt.Printf("  %v\n", swatch.HSL())
			if i == len(swatches)-1 {
				fmt.Println("")
			}
		}
	}
	fmt.Printf("%s: %06s\n", formatTarget("Vibrant", vibrant.Vibrant), formatPackgedRGB(palette.VibrantColor(0)))
	if palette.VibrantSwatch() != nil {
		fmt.Printf("  %v\n", palette.VibrantSwatch().HSL())
	}
	fmt.Printf("%s: %06s\n", formatTarget("Light Vibrant", vibrant.LightVibrant), formatPackgedRGB(palette.LightVibrantColor(0)))
	if palette.LightVibrantSwatch() != nil {
		fmt.Printf("  %v\n", palette.LightVibrantSwatch().HSL())
	}
	fmt.Printf("%s: %06s\n", formatTarget("Dark Vibrant", vibrant.DarkVibrant), formatPackgedRGB(palette.DarkVibrantColor(0)))
	if palette.DarkVibrantSwatch() != nil {
		fmt.Printf("  %v\n", palette.DarkVibrantSwatch().HSL())
	}
	fmt.Printf("%s: %06s\n", formatTarget("Muted", vibrant.Muted), formatPackgedRGB(palette.MutedColor(0)))
	if palette.MutedSwatch() != nil {
		fmt.Printf("  %v\n", palette.MutedSwatch().HSL())
	}
	fmt.Printf("%s: %06s\n", formatTarget("Light Muted", vibrant.LightMuted), formatPackgedRGB(palette.LightMutedColor(0)))
	if palette.LightMutedSwatch() != nil {
		fmt.Printf("  %v\n", palette.LightMutedSwatch().HSL())
	}
	fmt.Printf("%s: %06s\n", formatTarget("Dark Muted", vibrant.DarkMuted), formatPackgedRGB(palette.DarkMutedColor(0)))
	if palette.DarkMutedSwatch() != nil {
		fmt.Printf("  %v\n", palette.DarkMutedSwatch().HSL())
	}

	if len(outputFile) > 0 {
		if resizeImageArea > 0 {
			inputImage = vibrant.ScaleImageDown(inputImage, resizeImageArea, scalerByName[scalerName])
		}

		outputImageRectangle := inputImage.Bounds()
		var outputImage draw.Image = image.NewPaletted(outputImageRectangle, colorPalette)

		draw.Draw(outputImage, outputImage.Bounds(), inputImage, image.ZP, draw.Src)

		if debug {
			maxPoint := outputImageRectangle.Max
			var palettedStartPoint image.Point

			if maxPoint.X > maxPoint.Y {
				// Landscape, display images stacked.
				palettedStartPoint = image.Pt(0, maxPoint.Y+1)
				maxPoint.Y *= 2
				maxPoint.Y++
			} else {
				// Square or Portrait, display images side by side.
				palettedStartPoint = image.Pt(maxPoint.X+1, 0)
				maxPoint.X *= 2
				maxPoint.X++
			}

			swatchesPerRow := int(maximumColorCount)
			swatchRows := 1
			targetSwatchDimension := int(math.Max(math.Ceil(float64(maxPoint.X)*0.04), 8))
			actualSwatchDimension := maxPoint.X / swatchesPerRow
			for actualSwatchDimension < targetSwatchDimension {
				if swatchesPerRow%2 == 0 {
					swatchesPerRow = swatchesPerRow / 2
				} else {
					swatchesPerRow = (swatchesPerRow + 1) / 2
				}
				actualSwatchDimension = maxPoint.X / swatchesPerRow
				swatchRows = swatchRows << 1
			}

			swatchStartY := maxPoint.Y + 1
			maxPoint.Y = maxPoint.Y + 2 + (actualSwatchDimension * (swatchRows + len(palette.Targets())))

			rgbaImage := image.NewRGBA(image.Rectangle{image.ZP, maxPoint})
			// Add background...
			draw.Draw(rgbaImage, rgbaImage.Bounds(), image.NewUniform(image.Black), image.ZP, draw.Src)
			// Add original...
			draw.Draw(rgbaImage, inputImage.Bounds(), inputImage, image.ZP, draw.Src)
			// Add paletted...
			draw.Draw(rgbaImage, outputImage.Bounds().Add(palettedStartPoint), outputImage, image.ZP, draw.Src)
			// Add swatches...
			swatchRectangle := image.Rect(0, 0, actualSwatchDimension, actualSwatchDimension)
			sort.Sort(hueSwatchSorter(swatches))
			for i, swatch := range swatches {
				swatchX := (i % swatchesPerRow) * actualSwatchDimension
				swatchY := (i / swatchesPerRow) * actualSwatchDimension
				draw.Draw(rgbaImage, swatchRectangle.Add(image.Pt(swatchX, swatchStartY+swatchY)), image.NewUniform(swatch.Color()), image.ZP, draw.Src)
			}
			// Add targets...
			targetStartY := maxPoint.Y - (actualSwatchDimension * len(palette.Targets()))
			targetRectangle := image.Rect(0, 0, maxPoint.X, actualSwatchDimension)
			for i, target := range []*vibrant.Target{vibrant.Vibrant, vibrant.LightVibrant, vibrant.DarkVibrant, vibrant.Muted, vibrant.LightMuted, vibrant.DarkMuted} {
				swatch := palette.SwatchForTarget(target)
				var targetColor color.Color
				if swatch == nil {
					targetColor = color.Black
				} else {
					targetColor = swatch.Color()
				}
				draw.Draw(rgbaImage, targetRectangle.Add(image.Pt(0, targetStartY+(i*actualSwatchDimension))), image.NewUniform(targetColor), image.ZP, draw.Src)
			}

			outputImage = rgbaImage
		}

		output, err := os.Create(outputFile)
		if os.IsExist(err) {
			fmt.Printf("%s already exists\n", inputFile)
			os.Exit(1)
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer output.Close()

		if strings.HasSuffix(strings.ToLower(outputFile), ".png") {
			err = png.Encode(output, outputImage)
		} else if strings.HasSuffix(strings.ToLower(outputFile), ".jpg") {
			err = jpeg.Encode(output, outputImage, &jpeg.Options{90})
		} else {
			err = fmt.Errorf("Unable to find encoder for output file: %s", outputFile)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func formatTarget(name string, target *vibrant.Target) string {
	return fmt.Sprintf("%s (%0.1f, %0.1f)", name, target.TargetSaturation(), target.TargetLightness())
}

func formatPackgedRGB(c uint32) string {
	return fmt.Sprintf("0x%06s", strings.ToUpper(strconv.FormatUint(uint64(c), 16)))
}

type populationSwatchSorter []*vibrant.Swatch

func (p populationSwatchSorter) Len() int           { return len(p) }
func (p populationSwatchSorter) Less(i, j int) bool { return p[i].Population() > p[j].Population() }
func (p populationSwatchSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type hueSwatchSorter []*vibrant.Swatch

func (p hueSwatchSorter) Len() int           { return len(p) }
func (p hueSwatchSorter) Less(i, j int) bool { return p[i].HSL().H < p[j].HSL().H }
func (p hueSwatchSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
