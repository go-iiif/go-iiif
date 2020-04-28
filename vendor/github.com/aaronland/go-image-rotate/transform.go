package rotate

import (
	"context"
	"github.com/aaronland/go-image-transform"
	"image"
	"net/url"
)

type RotateTransformation struct {
	transform.Transformation
	orientation string
}

func init() {

	ctx := context.Background()
	err := transform.RegisterTransformation(ctx, "rotate", NewRotateTransformation)

	if err != nil {
		panic(err)
	}
}

func NewRotateTransformation(ctx context.Context, str_url string) (transform.Transformation, error) {

	parsed, err := url.Parse(str_url)

	if err != nil {
		return nil, err
	}

	query := parsed.Query()
	orientation := query.Get("orientation")

	if orientation != "" {
		orientation = "1"
	}

	tr := &RotateTransformation{
		orientation: orientation,
	}

	return tr, nil
}

func (tr *RotateTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return RotateImageWithOrientation(ctx, im, tr.orientation)
}
