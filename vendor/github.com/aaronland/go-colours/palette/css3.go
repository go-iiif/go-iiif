package palette

import (
	"context"
)

func init() {
	ctx := context.Background()
	err := RegisterPalette(ctx, "css3", NewCSS3Palette)
	if err != nil {
		panic(err)
	}
}

func NewCSS3Palette(ctx context.Context, uri string) (Palette, error) {

	data, err := FS.ReadFile("css3.json")

	if err != nil {
		return nil, err
	}

	return NewCommonPalette(data)
}
