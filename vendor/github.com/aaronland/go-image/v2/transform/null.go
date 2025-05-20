package transform

import (
	"context"
	"image"
)

func init() {
	ctx := context.Background()
	RegisterTransformation(ctx, "null", NewNullTransformation)
}

// NullTransformation is a struct that implements the `Transformation` interface that
// does not apply any transformations to images.
type NullTransformation struct {
	Transformation
}

// NewNullWriter returns a new `NullTransformation` instance.
// 'uri' in the form of:
//
//	null://
func NewNullTransformation(ctx context.Context, uri string) (Transformation, error) {

	tr := &NullTransformation{}
	return tr, nil
}

// Tranform returns 'im'.
func (tr *NullTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return im, nil
}
