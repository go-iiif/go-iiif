package colours

import (
	"context"
	"net/url"
)

func NewHexColour(ctx context.Context, hex string) (Colour, error) {

	q := url.Values{}
	q.Set("hex", hex)

	u := url.URL{}
	u.Scheme = "common"
	u.RawQuery = q.Encode()

	return NewColour(ctx, u.String())
}
