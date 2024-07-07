package rotate

import (
	"context"
	"fmt"
	"image"
	"net/url"

	"github.com/aaronland/go-image/transform"
)

// RotateTransformation is a struct that implements the `Transformation` interface for
// rotating images.
type RotateTransformation struct {
	transform.Transformation
	orientation string
}

func init() {
	ctx := context.Background()
	transform.RegisterTransformation(ctx, "rotate", NewRotateTransformation)
}

// NewRotateWriter returns a new `RotateTransformation` instance configure by 'uri'
// in the form of:
//
//	rotate://?orientation={ORIENTATION}
//
// Where {ORIENTATION} is expected to be a valid EXIF orientation string (1-8).
func NewRotateTransformation(ctx context.Context, uri string) (transform.Transformation, error) {

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	query := parsed.Query()
	orientation := query.Get("orientation")

	if orientation == "" {
		orientation = "1"
	}

	tr := &RotateTransformation{
		orientation: orientation,
	}

	return tr, nil
}

// Transform will rotate 'im' and return a new `image.Image` instance.
func (tr *RotateTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return RotateImageWithOrientation(ctx, im, tr.orientation)
}
