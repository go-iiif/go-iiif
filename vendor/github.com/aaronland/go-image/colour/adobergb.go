package colour

import (
	"context"
	"image"
	"runtime"

	"github.com/aaronland/go-image/transform"
	"github.com/mandykoh/prism"
	"github.com/mandykoh/prism/adobergb"
	"github.com/mandykoh/prism/srgb"
)

func init() {
	ctx := context.Background()
	transform.RegisterTransformation(ctx, "adobergb", NewAdobeRGBTransformation)
}

// AdobeRGBTransformation implements the `transform.Transformation` interface for converting
// all the colours in an image to match the Adobe RGB colour profile.
type AdobeRGBTransformation struct {
	transform.Transformation
}

// NewAdobeRGBTransformation returns a new `AdobeRGBTransformationTransformation` instance
// configured by 'uri' which is expected to take the form of:
//
//	adobergb://
func NewAdobeRGBTransformation(ctx context.Context, uri string) (transform.Transformation, error) {
	tr := &AdobeRGBTransformation{}
	return tr, nil
}

// Transform converts all the colours in 'im' to match the Adobe RGB colour profile.
func (tr *AdobeRGBTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	new_im := ToAdobeRGB(im)
	return new_im, nil
}

// ToAdobeRGB converts all the colours in 'im' to match the Adobe RGB colour profile.
func ToAdobeRGB(im image.Image) image.Image {

	input_im := prism.ConvertImageToNRGBA(im, runtime.NumCPU())
	new_im := image.NewNRGBA(input_im.Rect)

	for i := input_im.Rect.Min.Y; i < input_im.Rect.Max.Y; i++ {

		for j := input_im.Rect.Min.X; j < input_im.Rect.Max.X; j++ {

			inCol, alpha := adobergb.ColorFromNRGBA(input_im.NRGBAAt(j, i))
			outCol := srgb.ColorFromXYZ(inCol.ToXYZ())
			new_im.SetNRGBA(j, i, outCol.ToNRGBA(alpha))
		}
	}

	return new_im
}
