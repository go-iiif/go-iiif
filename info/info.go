package info

import (
	"fmt"
	iiifimage "github.com/go-iiif/go-iiif/v5/image"
	iiiflevel "github.com/go-iiif/go-iiif/v5/level"
	iiifprofile "github.com/go-iiif/go-iiif/v5/profile"
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
	Height   int                    `json:"height"`
	Width    int                    `json:"width"`
	Context  string                 `json:"@context"`
	Id       string                 `json:"@id"`
	Protocol string                 `json:"protocol"`
	Profiles []*iiifprofile.Profile `json:"profile"`
	Tiles    []*Tile                `json:"tiles"`
	Sizes    []*Size                `json:"sizes"`
}

func New(l iiiflevel.Level, im iiifimage.Image) (*Info, error) {

	pr, err := l.Profile()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new profile for level, %w", err)
	}

	dims, err := im.Dimensions()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive dimensions for image, %w", err)
	}

	i := &Info{
		Context:  IMAGE_V2_CONTEXT,
		Protocol: IMAGE_PROTOCOL,
		Id:       filepath.Join(l.Endpoint(), im.Identifier()),
		Profiles: []*iiifprofile.Profile{pr},
		Height:   dims.Height(),
		Width:    dims.Width(),
	}

	return i, nil
}
