package image

import (
       "github.com/fogleman/primitive/primitive"
       "log"
       "math"
       "time"
)

func PrimitiveImage(im Image) error {

        dims, err := im.Dimensions()

	if err != nil {
		return err
	}

	goimg, err := IIIFImageToGolangImage(im)

	if err != nil {
		return err
	}

	// Please make me config variables

	Alpha := 128
	Mode := 4	// mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect
	Number := 100

	h := float64(dims.Height())
	w := float64(dims.Width())
	max := math.Max(h, w)

	OutputSize := int(max)

	// See this - we're not dealing with animations yet

	t1 := time.Now()
	log.Println("starting model at", t1)

        model := primitive.NewModel(goimg, Alpha, OutputSize, primitive.Mode(Mode))

	for i := 1; i <= Number; i++ {

		tx := time.Since(t1)
		log.Printf("finished step %d in %v\n", i, tx)

		model.Step()
	}

	t2 := time.Since(t1)
	log.Println("finished model in", t2)

	return GolangImageToIIIFImage(model.Context.Image(), im)
}
