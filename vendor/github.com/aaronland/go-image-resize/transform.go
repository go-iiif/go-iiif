package resize

import (
	"context"
	"errors"
	"github.com/aaronland/go-image-transform"
	"image"
	"net/url"
	"strconv"
)

type ResizeTransformation struct {
	transform.Transformation
	max int
}

func init() {

	ctx := context.Background()
	err := transform.RegisterTransformation(ctx, "resize", NewResizeTransformation)

	if err != nil {
		panic(err)
	}
}

func NewResizeTransformation(ctx context.Context, str_url string) (transform.Transformation, error) {

	parsed, err := url.Parse(str_url)

	if err != nil {
		return nil, err
	}

	query := parsed.Query()
	str_max := query.Get("max")

	if str_max == "" {
		return nil, errors.New("Missing parameter: max")
	}

	max, err := strconv.Atoi(str_max)

	if err != nil {
		return nil, err
	}

	tr := &ResizeTransformation{
		max: max,
	}

	return tr, nil
}

func (tr *ResizeTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return ResizeImageMax(ctx, im, tr.max)
}
