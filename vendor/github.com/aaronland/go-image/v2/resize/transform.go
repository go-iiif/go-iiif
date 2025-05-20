package resize

import (
	"context"
	"fmt"
	"image"
	"net/url"
	"strconv"

	"github.com/aaronland/go-image/v2/transform"
)

// ResizeTransformation is a struct that implements the `Transformation` interface for
// resizing images by a maximum dimension.
type ResizeTransformation struct {
	transform.Transformation
	max int
}

func init() {

	ctx := context.Background()
	transform.RegisterTransformation(ctx, "resize", NewResizeTransformation)
}

// NewResizeWriter returns a new `ResizeTransformation` instance configure by 'uri'
// in the form of:
//
//	resize://?max={MAX}
//
// Where {MAX} is expected to be a maximum dimension for the resized image.
func NewResizeTransformation(ctx context.Context, str_url string) (transform.Transformation, error) {

	parsed, err := url.Parse(str_url)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	query := parsed.Query()
	str_max := query.Get("max")

	if str_max == "" {
		return nil, fmt.Errorf("Missing parameter: max")
	}

	max, err := strconv.Atoi(str_max)

	if err != nil {
		return nil, fmt.Errorf("Failed to convert ?max= parameter, %w", err)
	}

	tr := &ResizeTransformation{
		max: max,
	}

	return tr, nil
}

// Transform will resize 'im' and return a new `image.Image` instance.
func (tr *ResizeTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return ResizeImage(ctx, im, tr.max)
}
