package palette

import (
	"encoding/json"

	"github.com/aaronland/go-colours"
)

type CommonPalette struct {
	Palette         `json:",omitempty"`
	CommonReference string                  `json:"reference"`
	CommonColours   []*colours.CommonColour `json:"colours,omitempty"`
}

func NewCommonPalette(data []byte) (Palette, error) {

	var p CommonPalette

	err := json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *CommonPalette) Reference() string {
	return p.CommonReference
}

func (p *CommonPalette) Colours() []colours.Colour {

	c := make([]colours.Colour, 0)

	for _, pc := range p.CommonColours {
		c = append(c, pc)
	}

	return c
}
