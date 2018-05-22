package atkinson

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func decodeImage(filePath string) (img image.Image, err error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	last3Strings := strings.ToLower(filePath[len(filePath)-3:])
	last4Strings := strings.ToLower(filePath[len(filePath)-4:])
	if last3Strings == "jpg" || last4Strings == "jpeg" {
		img, err = jpeg.Decode(reader)
	} else if last3Strings == "gif" {
		img, err = gif.Decode(reader)
	} else if last3Strings == "png" {
		img, err = png.Decode(reader)
	} else {
		img = nil
		err = errors.New("unknown format")
	}
	return
}

func DitherFile(path string) (result image.Image, err error) {
	img, err := decodeImage(path)
	if err != nil {
		return nil, err
	}

	return Dither(img)
}

func Dither(img image.Image) (result image.Image, err error) {
	bounds := img.Bounds()
	dstImg := image.NewGray(bounds)
	draw.Draw(dstImg, bounds, img, image.ZP, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := dstImg.At(x, y)
			gray, ok := c.(color.Gray)
			if ok != true {
				continue
			}

			var m = mono(gray.Y)
			quant_err := int((float64(gray.Y) - float64(m)) / 8)
			dstImg.SetGray(x, y, color.Gray{m})

			s := 1
			neighborsPoint := []image.Point{
				image.Point{x + s, y},
				image.Point{x - s, y + s},
				image.Point{x, y + s},
				image.Point{x + s, y + s},
				image.Point{x + 2*s, y},
				image.Point{x, y + 2*s},
			}
			for _, p := range neighborsPoint {
				neighborColor := dstImg.At(p.X, p.Y)
				gray, ok := neighborColor.(color.Gray)
				if ok != true {
					continue
				}
				dstImg.SetGray(p.X, p.Y, color.Gray{uint8(int16(gray.Y) + int16(quant_err))})
			}
		}
	}

	return dstImg, nil
}

//func luminance(r, g, b uint8) uint8 {
//	return uint8(float64(r)*0.3 + float64(g)*0.59 + float64(b)*0.11)
//}

func mono(l uint8) uint8 {
	if l < 128 {
		return 0
	} else {
		return 255
	}
}
