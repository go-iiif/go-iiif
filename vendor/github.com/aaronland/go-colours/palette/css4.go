package palette

import (
	"context"
)

func init() {
	ctx := context.Background()
	err := RegisterPalette(ctx, "css4", NewCSS4Palette)
	if err != nil {
		panic(err)
	}
}

func NewCSS4Palette(ctx context.Context, uri string) (Palette, error) {

	data, err := FS.ReadFile("css4.json")

	if err != nil {
		return nil, err
	}

	return NewCommonPalette(data)
}
