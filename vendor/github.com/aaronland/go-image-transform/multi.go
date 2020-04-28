package transform

import (
	"context"
	"image"
	_ "log"
)

type MultiTransformation struct {
	Transformation
	transforms []Transformation
}

func NewMultiTransformation(ctx context.Context, transforms ...Transformation) (Transformation, error) {

	tr := &MultiTransformation{
		transforms: transforms,
	}

	return tr, nil
}

func (tr *MultiTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {

	var err error

	for _, t := range tr.transforms {

		im, err = t.Transform(ctx, im)

		if err != nil {
			return nil, err
		}
	}

	return im, nil
}
