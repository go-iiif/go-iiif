package colours

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

type CommonColour struct {
	Colour          `json:",omitempty"`
	CommonName      string   `json:"name,omitempty"`
	CommonHex       string   `json:"hex"`
	CommonReference string   `json:"reference,omitempty"`
	CommonClosest   []Colour `json:"closest,omitempty"`
}

func init() {
	ctx := context.Background()
	err := RegisterColour(ctx, "common", NewCommonColour)
	if err != nil {
		panic(err)
	}
}

func NewCommonColour(ctx context.Context, uri string) (Colour, error) {

	hex := "###"
	name := "unknown"
	ref := "unknown"

	var closest []Colour

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	if q.Has("hex") {
		hex = q.Get("hex")
	}

	if q.Has("name") {
		name = q.Get("name")
	}

	if q.Has("ref") {
		ref = q.Get("ref")
	}

	closest = make([]Colour, 0)

	/*

		for i, c := range args[3:] {
			closest[i] = c.(Colour)
		}
	*/

	if !strings.HasPrefix(hex, "#") {
		hex = fmt.Sprintf("#%s", hex)
	}

	_, err = colorful.Hex(hex)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive hex (%s), %w", hex, err)
	}

	hc := CommonColour{
		CommonName:      name,
		CommonHex:       hex,
		CommonReference: ref,
		CommonClosest:   closest,
	}

	return &hc, nil
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

/*

 */
