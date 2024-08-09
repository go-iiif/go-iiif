package transform

import (
	"context"
	"fmt"
	"image"
)

// Type MultiTransformation implements the `Transformation` interface for applying multiple transformations to images.
type MultiTransformation struct {
	Transformation
	transforms []Transformation
}

// NewMultiTransformationWithURIs returns a `MultiTransformation` instance derived from 'transformation_uris'.
// Transformations are applied in the same order as URIs defined in 'transformation_uris'.
func NewMultiTransformationWithURIs(ctx context.Context, transformation_uris ...string) (Transformation, error) {

	transformations := make([]Transformation, len(transformation_uris))

	for idx, uri := range transformation_uris {

		t, err := NewTransformation(ctx, uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create transformation for %s, %v", uri, err)
		}

		transformations[idx] = t
	}

	return NewMultiTransformation(ctx, transformations...)
}

// NewMultiTransformationWithURIs returns a `MultiTransformation` instance derived from 'transformation'.
// Transformations are applied in the same order as instances defined in 'transformation'.
func NewMultiTransformation(ctx context.Context, transformations ...Transformation) (Transformation, error) {

	tr := &MultiTransformation{
		transforms: transformations,
	}

	return tr, nil
}

func (tr *MultiTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {

	var err error

	for idx, t := range tr.transforms {

		im, err = t.Transform(ctx, im)

		if err != nil {
			return nil, fmt.Errorf("Failed to apply transform at offset %d, %w", idx, err)
		}
	}

	return im, nil
}
