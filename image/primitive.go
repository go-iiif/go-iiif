package image

// "primitive": { "syntax": "primitive:mode,iterations,alpha", "required": false, "supported": true, "match": "^primitive\:[0-4]\,\d+,\d+$" }
// mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect

import (
	"github.com/fogleman/primitive/primitive"
	"log"
	"math"
	"time"
)

type PrimitiveOptions struct {
	Alpha      int
	Mode       int
	Iterations int
	Size       int
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

	// See this - we're not dealing with animations yet

	t1 := time.Now()
	log.Println("starting model at", t1)

	model := primitive.NewModel(goimg, alpha, size, primitive.Mode(mode))

	for i := 1; i <= opts.Iterations; i++ {

		tx := time.Since(t1)
		log.Printf("finished step %d in %v\n", i, tx)

		model.Step()
	}

	t2 := time.Since(t1)
	log.Println("finished model in", t2)

	return GolangImageToIIIFImage(model.Context.Image(), im)
}
