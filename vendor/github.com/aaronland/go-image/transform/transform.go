// Package transform provides interfaces for applying transformations to images.
package transform

import (
	"context"
	"fmt"
	"image"
	"net/url"

	"github.com/aaronland/go-roster"
)

var transformations roster.Roster

// TransformationInitializationFunc is a function defined by individual transformation package and used to create
// an instance of that transformation
type InitializeTransformationFunc func(context.Context, string) (Transformation, error)

// Transformation is an interface for writing data to multiple sources or targets.
type Transformation interface {
	// Transform applies a transformation to an `image.Image` instance returning a new `image.Image` instance.
	Transform(context.Context, image.Image) (image.Image, error)
}

func ensureRoster() error {

	if transformations == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new transformations roster, %w", err)
		}

		transformations = r
	}

	return nil
}

// RegisterTransformation registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Transformation` instances by the `NewTransformation` method.
func RegisterTransformation(ctx context.Context, name string, f InitializeTransformationFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return transformations.Register(ctx, name, f)
}

// NewTransformation returns a new `Transformation` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `TransformationInitializationFunc`
// function used to instantiate the new `Transformation`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterTransformation` method.
func NewTransformation(ctx context.Context, uri string) (Transformation, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	scheme := u.Scheme

	i, err := transformations.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive transformation for %s, %w", scheme, err)
	}

	if i == nil {
		return nil, fmt.Errorf("Undefined transformation for %s", scheme)
	}

	f := i.(InitializeTransformationFunc)
	return f(ctx, uri)
}
