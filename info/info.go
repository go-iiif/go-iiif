package info

import (
	"encoding/json"
	"fmt"
	iiifimage "github.com/go-iiif/go-iiif/v5/image"
	iiiflevel "github.com/go-iiif/go-iiif/v5/level"
	iiifservice "github.com/go-iiif/go-iiif/v5/service"
	"github.com/tidwall/pretty"
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
	Height   int                   `json:"height"`
	Width    int                   `json:"width"`
	Context  string                `json:"@context"`
	Id       string                `json:"@id"`
	Protocol string                `json:"protocol"`
	Profile  []interface{}         `json:"profile"`
	Tiles    []*Tile               `json:"tiles,omitempty"`
	Sizes    []*Size               `json:"sizes,omitempty"`
	Services []iiifservice.Service `json:"service,omitempty"`
}

func New(iiif_context string, l iiiflevel.Level, im iiifimage.Image) (*Info, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive dimensions for image, %w", err)
	}

	i := &Info{
		Context:  iiif_context,
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

func  MarshalJSON(i *Info) ([]byte, error) {

	body, err := json.Marshal(i)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal info, %w", err)
	}

	body = pretty.Pretty(body)
	return body, nil
}
