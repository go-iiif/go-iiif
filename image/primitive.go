package image

// "primitive": { "syntax": "primitive:mode,iterations,alpha", "required": false, "supported": true, "match": "^primitive\:[0-4]\,\d+,\d+$" }
// mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect

import (
	"bytes"
	"github.com/fogleman/primitive/primitive"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "log"
	"math"
	"runtime"
	_ "time"
)

type PrimitiveOptions struct {
	Alpha      int
	Mode       int
	Iterations int
	Size       int
	Animated   bool
}

func PrimitiveImage(im Image, opts PrimitiveOptions) error {

	dims, err := im.Dimensions()

	if err != nil {
		return err
	}

	goimg, err := IIIFImageToGolangImage(im)

	if err != nil {
		return err
	}

	alpha := opts.Alpha
	mode := opts.Mode
	size := opts.Size

	if size == 0 {
		h := float64(dims.Height())
		w := float64(dims.Width())
		max := math.Max(h, w)
		size = int(max)
	}

	// t1 := time.Now()
	// log.Println("starting model at", t1)

	workers := runtime.NumCPU()

	bg := primitive.MakeColor(primitive.AverageImageColor(goimg))
	model := primitive.NewModel(goimg, bg, size, workers)

	for i := 1; i <= opts.Iterations; i++ {

		// tx := time.Since(t1)
		// log.Printf("finished step %d in %v\n", i, tx)

		model.Step(primitive.ShapeType(mode), alpha, workers)
	}

	// t2 := time.Since(t1)
	// log.Println("finished model in", t2)

	if opts.Animated {

		g := gif.GIF{}

		frames := model.Frames(0.001)

		delay := 25
		lastDelay := delay * 10

		for i, src := range frames {

			// the original code in primitive/utils.go
			// dst := image.NewPaletted(src.Bounds(), palette.Plan9)
			// draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)

			// https://groups.google.com/forum/#!topic/golang-nuts/28Kk1FfG5XE
			// https://github.com/golang/go/blob/master/src/image/gif/writer.go#L358-L366

			opts := gif.Options{
				NumColors: 256,
				Drawer:    draw.FloydSteinberg,
				Quantizer: nil,
			}

			dst := image.NewPaletted(src.Bounds(), palette.Plan9[:opts.NumColors])
			opts.Drawer.Draw(dst, dst.Rect, src, image.ZP)

			g.Image = append(g.Image, dst)

			if i == len(frames)-1 {
				g.Delay = append(g.Delay, lastDelay)
			} else {
				g.Delay = append(g.Delay, delay)
			}
		}

		out := new(bytes.Buffer)
		err := gif.EncodeAll(out, &g)

		if err != nil {
			return err
		}

		return im.Update(out.Bytes())

	} else {
		goimg := model.Context.Image()
		return GolangImageToIIIFImage(goimg, im)
	}

}
