package level

import (
	"fmt"
	iiifcompliance "github.com/thisisaaronland/go-iiif/compliance"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"log"
)

type Level2 struct {
	Level      `json:"omit"`
	Context    string                    `json:@profile`
	Id         string                    `json:"@id"`
	Type       string                    `json:"@type"` // Optional or iiif:Image
	Formats    []string                  `json:"formats"`
	Qualities  []string                  `json:"qualities"`
	Supports   []string                  `json:"supports"`
	compliance iiifcompliance.Compliance `json:"omit"`
}

/*
var re_alpha *regexp.Regexp
var re_region *regexp.Regexp
var re_size *regexp.Regexp
var re_rotation *regexp.Regexp
var re_quality *regexp.Regexp

func init() {

	var err error

	re_alpha, err = regexp.Compile(`^[a-z]+$`)

	if err != nil {
		log.Fatal(err)
	}

	re_region, err = regexp.Compile(`^(?:full|square|\d+\,\d+\,\d+\,\d+|pct\:\d+(\.\d+)?,\d+(\.\d+)?,\d+(\.\d+)?,\d+(\.\d+)?)$`)

	if err != nil {
		log.Fatal(err)
	}

	re_rotation, err = regexp.Compile(`^\!?\d+`)

	if err != nil {
		log.Fatal(err)
	}

	re_size, err = regexp.Compile(`^(?:(?:max|full)|(?:\d+\,\d+)|(?:\!\d+\,\d+)|(\d+\,)|(\,\d+)|(pct\:\d+))$`)

	if err != nil {
		log.Fatal(err)
	}

}

*/

func NewLevel2(config *iiifconfig.Config, host string) (*Level2, error) {

	compliance, err := iiifcompliance.NewLevel2Compliance(config)

	if err != nil {
		log.Fatal(err)
	}

	id := fmt.Sprintf("http://%s/level2.json", host)

	l := Level2{
		Context:    "http://iiif.io/api/image/2/context.json",
		Id:         id,
		Type:       "iiif:ImageProfile",
		Formats:    compliance.Formats(),
		Qualities:  compliance.Qualities(),
		Supports:   []string{},
		compliance: compliance,
	}

	return &l, nil
}

func (l *Level2) Compliance() iiifcompliance.Compliance {
	return l.compliance
}

/*

// full
// square
// x,y,w,h (in pixels)
// pct:x,y,w,h (in percents)

func (l *Level2) IsValidImageRegion(region string) (bool, error) {

	if !re_region.MatchString(region) {
		message := fmt.Sprintf("Invalid IIIF 2.1 region: %s", region)
		return false, errors.New(message)
	}

	return true, nil
}

// max, full
// w,h (deform)
// !w,h (best fit within size)
// w, (force width)
// ,h (force height)
// pct:n (resize)

func (l *Level2) IsValidImageSize(size string) (bool, error) {

	if !re_size.MatchString(size) {
		message := fmt.Sprintf("Invalid IIIF 2.1 size: %s", size)
		return false, errors.New(message)
	}

	return true, nil
}

// n angle clockwise in degrees
// !n angle clockwise in degrees with a flip (beforehand)

func (l *Level2) IsValidImageRotation(rotation string) (bool, error) {

	if !re_rotation.MatchString(rotation) {
		message := fmt.Sprintf("Invalid IIIF 2.1 rotation: %s", rotation)
		return false, errors.New(message)
	}

	parsed, err := strconv.ParseInt(strings.Trim(rotation, "!"), 10, 64)

	if err != nil {
		return false, err
	}

	if parsed > 360 {
		message := fmt.Sprintf("Invalid IIIF 2.1 rotation: %s", rotation)
		return false, errors.New(message)
	}

	return true, nil
}

func (l *Level2) IsValidImageQuality(quality string) (bool, error) {

	if !re_alpha.MatchString(quality) {
		message := fmt.Sprintf("Invalid IIIF 2.1 quality: %s", quality)
		return false, errors.New(message)
	}

	ok := false

	for _, test := range l.Qualities {

		if quality == test {
			ok = true
			break
		}
	}

	if !ok {
		message := fmt.Sprintf("Unsupported IIIF 2.1 quality: %s", quality)
		return false, errors.New(message)
	}

	return true, nil
}

func (l *Level2) IsValidImageFormat(format string) (bool, error) {

	if !re_alpha.MatchString(format) {
		message := fmt.Sprintf("Invalid IIIF 2.1 format: %s", format)
		return false, errors.New(message)
	}

	ok := false

	for _, test := range l.Formats {

		if format == test {
			ok = true
			break
		}
	}

	if !ok {
		message := fmt.Sprintf("Unsupported IIIF 2.1 format: %s", format)
		return false, errors.New(message)
	}

	return true, nil
}

*/
