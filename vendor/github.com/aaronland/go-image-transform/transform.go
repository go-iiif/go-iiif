package transform

import (
	"context"
	"github.com/aaronland/go-roster"
	"image"
	"net/url"
)

type InitializeTransformationFunc func(context.Context, string) (Transformation, error)

type Transformation interface {
	Transform(context.Context, image.Image) (image.Image, error)
}

var transformations roster.Roster

func ensureRoster() error {

	if transformations == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		transformations = r
	}

	return nil
}

func RegisterTransformation(ctx context.Context, name string, f InitializeTransformationFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return transformations.Register(ctx, name, f)
}

func NewTransformation(ctx context.Context, uri string) (Transformation, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := transformations.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(InitializeTransformationFunc)
	return f(ctx, uri)
}
