package palette

import (
	"image"
)

type Palette interface {
	Extract(image.Image) ([]Color, error)
}

type Color struct {
	Color   string `json:"color"`
	Closest string `json:"closest,omitempty"`
}
