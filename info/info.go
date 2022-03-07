package info

import (
	"fmt"
	iiifimage "github.com/go-iiif/go-iiif/v5/image"
	iiiflevel "github.com/go-iiif/go-iiif/v5/level"
	"path/filepath"
)

const IMAGE_PROTOCOL string = "http://iiif.io/api/image"

const IMAGE_V2_CONTEXT string = "http://iiif.io/api/image/2/context.json"

type Size struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type Tile struct {
	Height       int   `json:"height"`
	Width        int   `json:"width"`
	ScaleFactors []int `json:"scaleFactors"`
}

type Info struct {
	Height   int           `json:"height"`
	Width    int           `json:"width"`
	Context  string        `json:"@context"`
	Id       string        `json:"@id"`
	Protocol string        `json:"protocol"`
	Profile  []interface{} `json:"profile"`
	Tiles    []*Tile       `json:"tiles"`
	Sizes    []*Size       `json:"sizes"`
}

func New(l iiiflevel.Level, im iiifimage.Image) (*Info, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive dimensions for image, %w", err)
	}

	i := &Info{
		Context:  IMAGE_V2_CONTEXT,
		Protocol: IMAGE_PROTOCOL,
		Id:       filepath.Join(l.Endpoint(), im.Identifier()),
		Profile: []interface{}{
			l.Profile(),
			l,
		},
		Height: dims.Height(),
		Width:  dims.Width(),
	}

	return i, nil
}
