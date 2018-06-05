package colours

import (
	"errors"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"strings"
)

type Colour interface {
	Name() string
	Hex() string
	Reference() string
	Closest() []Colour
	AppendClosest(Colour) error // I don't love this... (20180605/thisisaaronland)
	String() string
}

type Palette interface {
	Reference() string
	Colours() []Colour
}

type Extruder interface {
	Colours(image.Image, int) ([]Colour, error)
}

type Grid interface {
	Closest(Colour, Palette) (Colour, error)
}

type CommonColour struct {
	Colour          `json:",omitempty"`
	CommonName      string   `json:"name,omitempty"`
	CommonHex       string   `json:"hex"`
	CommonReference string   `json:"reference,omitempty"`
	CommonClosest   []Colour `json:"closest,omitempty"`
}

func (hc *CommonColour) Name() string {
	return hc.CommonName
}

func (hc *CommonColour) Hex() string {
	return hc.CommonHex
}

func (hc *CommonColour) Reference() string {
	return hc.CommonReference
}

func (hc *CommonColour) AppendClosest(c Colour) error {

	if hc.CommonClosest == nil {
		hc.CommonClosest = make([]Colour, 0)
	}

	hc.CommonClosest = append(hc.CommonClosest, c)
	return nil
}

func (hc *CommonColour) String() string {

	name := hc.Name()
	hex := hc.Hex()
	ref := hc.Reference()

	if name == hex {
		return hex
	}

	return fmt.Sprintf("%s (%s, %s)", hex, name, ref)
}

func NewColour(args ...interface{}) (Colour, error) {

	hex := "###"
	name := "unknown"
	ref := "unknown"

	var closest []Colour

	var err error

	switch len(args) {

	case 0:
		err = errors.New("Insuffient arguments")
	case 1:
		hex = args[0].(string)
	case 2:
		hex = args[0].(string)
		name = args[1].(string)
	case 3:
		hex = args[0].(string)
		name = args[1].(string)
		ref = args[2].(string)
	default:
		hex = args[0].(string)
		name = args[1].(string)
		ref = args[2].(string)

		closest = make([]Colour, 0)

		for i, c := range args[3:] {
			closest[i] = c.(Colour)
		}
	}

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(hex, "#") {
		hex = fmt.Sprintf("#%s", hex)
	}

	_, err = colorful.Hex(hex)

	if err != nil {
		return nil, err
	}

	hc := CommonColour{
		CommonName:      name,
		CommonHex:       hex,
		CommonReference: ref,
		CommonClosest:   closest,
	}

	return &hc, nil
}

type CommonPalette struct {
	Palette         `json:",omitempty"`
	CommonReference string          `json:"reference"`
	CommonColours   []*CommonColour `json:"colours,omitempty"`
}

func (p *CommonPalette) Reference() string {
	return p.CommonReference
}

func (p *CommonPalette) Colours() []Colour {

	// Y DO I NEED TO DOOOOOOOOOOOOOOOOOOO THIS???
	// Y U SO WEIRD GOOOOOOOOOOOOOOOO????
	// (20180605/thisisaaronland)

	c := make([]Colour, 0)

	for _, pc := range p.CommonColours {
		c = append(c, pc)
	}

	return c
}
