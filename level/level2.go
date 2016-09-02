package level

import (
	"fmt"
)

type Level2 struct {
	Level
	Context   string   `json:@profile`
	Id        string   `json:"@id"`
	Type      string   `json:"@type"` // Optional or iiif:Image
	Formats   []string `json:"formats"`
	Qualities []string `json:"qualities"`
	Supports  []string `json:"supports"`
}

func NewLevel2(host string) (*Level2, error) {

	id := fmt.Sprintf("http://%s/level2.json", host)

	l := Level2{
		Context:   "http://iiif.io/api/image/2/context.json",
		Id:        id,
		Type:      "iiif:ImageProfile",
		Formats:   []string{"jpg", "png", "webp"},
		Qualities: []string{"gray", "default"},
		Supports:  []string{},
	}

	return &l, nil
}
