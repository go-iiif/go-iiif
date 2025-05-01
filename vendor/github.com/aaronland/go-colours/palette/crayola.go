package palette

import (
	"context"
)

func init() {
	ctx := context.Background()
	err := RegisterPalette(ctx, "crayola", NewCrayolaPalette)
	if err != nil {
		panic(err)
	}
}

func NewCrayolaPalette(ctx context.Context, uri string) (Palette, error) {

	data, err := FS.ReadFile("crayola.json")

	if err != nil {
		return nil, err
	}

	return NewCommonPalette(data)
}
