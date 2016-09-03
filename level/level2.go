package level

// http://iiif.io/api/image/2.1/compliance/

import (
	"errors"
	"fmt"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type ComplianceDetails struct {
	Name     string `json:"name"`
	Syntax   string `json:"syntax"`
	Required bool   `json:"required"`
	Match    string `json:"match,omitempty"`
}

type Level2ImageCompliance struct {
	Region   []ComplianceDetails `json:"region"`
	Size     []ComplianceDetails `json:"size"`
	Rotation []ComplianceDetails `json:"rotation"`
	Quality  []ComplianceDetails `json:"quality"`
	Format   []ComplianceDetails `json:"format"`
}

type Level2Compliance struct {
	Image Level2ImageCompliance `json:"image"`
}

type Level2 struct {
	Level
	Context   string   `json:@profile`
	Id        string   `json:"@id"`
	Type      string   `json:"@type"` // Optional or iiif:Image
	Formats   []string `json:"formats"`
	Qualities []string `json:"qualities"`
	Supports  []string `json:"supports"`
}

var compliance_spec = `{
    "image": {
    	     "region": [
	     	       { "name": "",             "syntax": "full",        "required": 1, "match": "^full$" },
		       { "name": "regionByPx",   "syntax": "x,y,w,h",     "required": 1, "match": "^\d+\,\d+\,\d+\,\d+$" },
		       { "name": "regionByPct",  "syntax": "pct:x,y,w,h", "required": 1, "match": "^pct\:\d+\,\d+\,\d+\,\d+$" },
		       { "name": "regionSquare", "syntax": "square",      "required": 0, "match": "^square$" }
	     },
	     "size": [
	     		{ "name": "",                  "syntax": "full",  "required": 1, "match": "^full$" },
	     		{ "name": "",                  "syntax": "max",   "required": 0, "match": "^max$" },
	     		{ "name": "sizeByW",           "syntax": "w,",    "required": 1, "match": "^\d+\,$" },			
	     		{ "name": "sizeByH",           "syntax": ",h",    "required": 1, "match": "^\,\d+$" },
	     		{ "name": "sizeByPct",         "syntax": "pct:n", "required": 1, "match": "^pct\:\d+(\.\d+)?$" },			
	     		{ "name": "sizeByConfinedWh",  "syntax": "!w,h",  "required": 1, "match": "" },
	     		{ "name": "sizeByDistortedWh", "syntax": "w,h",   "required": 1, "match": "" },			
	     		{ "name": "sizeByWh",          "syntax": "w,h",   "required": 1, "match": "" },
	     		{ "name": "sizeAboveFull",     "syntax": "",      "required": 0, "match": "" }
	     ],
	     "rotation": [
	     		{ "name": "",                  "syntax": "0",          "required": 1, "match": "" },
	     		{ "name": "rotationBy90s",     "syntax": "90,180,270", "required": 1, "match": "" },
	     		{ "name": "rotationArbitrary", "syntax": "",           "required": 0, "match": "" },			
	     		{ "name": "mirroring",         "syntax": "!n",         "required": 1, "match": "" }
	     ],
	     "quality": [
	     		{ "name": "", "syntax": "default", "required": 1, "match": "" },
	     		{ "name": "", "syntax": "color",   "required": 0, "match": "" },
	     		{ "name": "", "syntax": "gray",    "required": 0, "match": "" },			
	     		{ "name": "", "syntax": "bitonal", "required": 1, "match": "" }
             ],
	     "format": [
	     	       { "name": "", "syntax": "jpg",  "required": 1, "match": "" },
       	     	       { "name": "", "syntax": "png",  "required": 1, "match": "" },
       	     	       { "name": "", "syntax": "tif",  "required": 0, "match": "" },
      	     	       { "name": "", "syntax": "gif",  "required": 0, "match": "" },
       	     	       { "name": "", "syntax": "pdf",  "required": 0, "match": "" },
      	     	       { "name": "", "syntax": "jp2",  "required": 0, "match": "" },
       	     	       { "name": "", "syntax": "webp", "required": 0, "match": "" },
	     ]	     
    },
    "http": [
    	    { "feature": "base URI redirects", "name": "baseUriRedirect", "required": 1 },
	    { "feature": "CORS", "name": "cors", "required": 1 },
	    { "feature": "json-ld media type", "name": "jsonldMediaType", "required": 1 },
	    { "feature": "profile link header", "name": "profileLinkHeader", "required": 0 },
	    { "feature": "canonical link header", "name": "canonicalLinkHeader", "required": 0 }
    ]	    
}`

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

	/*
		re_quality, err = regexp.Compile(`^(?:color|grey|bitonal|default|dither)$`)

		if err != nil {
			log.Fatal(err)
		}
	*/

	re_size, err = regexp.Compile(`^(?:(?:max|full)|(?:\d+\,\d+)|(?:\!\d+\,\d+)|(\d+\,)|(\,\d+)|(pct\:\d+))$`)

	if err != nil {
		log.Fatal(err)
	}

}

func NewLevel2(config *iiifconfig.Config, host string) (*Level2, error) {

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
