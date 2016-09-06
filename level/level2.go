package level

// http://iiif.io/api/image/2.1/compliance/

import (
       "encoding/json"
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
	Supported bool  `json:"supported"`
	Match    string `json:"match,omitempty"`
}

type Level2ImageCompliance struct {
	Region   map[string]ComplianceDetails `json:"region"`
	Size     map[string]ComplianceDetails `json:"size"`
	Rotation map[string]ComplianceDetails `json:"rotation"`
	Quality  map[string]ComplianceDetails `json:"quality"`
	Format   map[string]ComplianceDetails `json:"format"`
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
    	     "region": {
	     	       "full": { "syntax": "full",        "required": true, "match": "^full$" },
		       "regionByPx": { "syntax": "x,y,w,h",     "required": true, "match": "^\\d+\\,\\d+\\,\\d+\\,\\d+$" },
		       "regionByPct": {  "syntax": "pct:x,y,w,h", "required": true, "match": "^pct\\:\\d+\\,\\d+\\,\\d+\\,\\d+$" },
		       "regionSquare": { "syntax": "square",      "required": false, "match": "^square$" }
	     },
	     "size": {
	     		"full": {              "syntax": "full",  "required": true, "supported": true, "match": "^full$" },
	     		"max": {               "syntax": "max",   "required": false, "supported": true, "match": "^max$" },
	     		"sizeByW": {           "syntax": "w,",    "required": true, "supported": true, "match": "^\\d+\\,$" },			
	     		"sizeByH": {           "syntax": ",h",    "required": true, "supported": true, "match": "^\\,\\d+$" },
	     		"sizeByPct": {         "syntax": "pct:n", "required": true, "supported": true, "match": "^pct\\:\\d+(\\.\\d+)?$" },			
	     		"sizeByConfinedWh": {  "syntax": "!w,h",  "required": true, "supported": true, "match": "" },
	     		"sizeByDistortedWh": { "syntax": "w,h",   "required": true, "supported": true, "match": "" },			
	     		"sizeByWh": {          "syntax": "w,h",   "required": true, "supported": true, "match": "" },
	     		"sizeAboveFull": {     "syntax": "",      "required": false, "supported": false, "match": "" }
	     },
	     "rotation": {
	     		"none": {              "syntax": "0",          "required": true, "supported": true, "match": "" },
	     		"rotationBy90s": {     "syntax": "90,180,270", "required": true, "supported": true, "match": "" },
	     		"rotationArbitrary": { "syntax": "",           "required": false, "supported": true, "match": "" },			
	     		"mirroring": {         "syntax": "!n",         "required": true, "supported": true, "match": "" }
	     },
	     "quality": {
	     		"default": { "syntax": "default", "required": true, "supported": true, "match": "^default$" },
	     		"color": { "syntax": "color",   "required": false, "supported": false, "match": "^colou?r$" },
	     		"gray": { "syntax": "gray",    "required": false, "supported": false, "match": "gr(?:e|a)y$" },			
	     		"bitonal" { "syntax": "bitonal", "required": true, "supported": true, "match": "^bitonal$" }
             },
	     "format": {
	     	       "jpg": { "syntax": "jpg",  "required": true, "supported": true, "match": "^jpe?g$" },
       	     	       "png": { "syntax": "png",  "required": true, "supported": true, "match": "^png$" },
       	     	       "tif": { "syntax": "tif",  "required": false, "supported": false, "match": "^tiff?$" },
      	     	       "gif": { "syntax": "gif",  "required": false, "supported": false, "match": "^gif$" },
       	     	       "pdf": { "syntax": "pdf",  "required": false, "supported": false, "match": "^pdf$" },
      	     	       "jp2": { "syntax": "jp2",  "required": false, "supported": false, "match": "^jp2$" },
       	     	       "webp": { "syntax": "webp", "required": false, "supported": false, "match": "^webp$" }
	     }	     
    }
    /*
    "http": [
    	    { "feature": "base URI redirects", "name": "baseUriRedirect", "required": true, "supported": 1 },
	    { "feature": "CORS", "name": "cors", "required": true, "supported": 1 },
	    { "feature": "json-ld media type", "name": "jsonldMediaType", "required": true, "supported": 1 },
	    { "feature": "profile link header", "name": "profileLinkHeader", "required": false, "supported": 0 },
	    { "feature": "canonical link header", "name": "canonicalLinkHeader", "required": false, "supported": 0 }
    ]
    */	    
}`

var re_alpha *regexp.Regexp
var re_region *regexp.Regexp
var re_size *regexp.Regexp
var re_rotation *regexp.Regexp
var re_quality *regexp.Regexp

func init() {

	var err error

     	c := Level2Compliance{}
	err = json.Unmarshal([]byte(compliance_spec), &c)

	if err != nil {
		log.Fatal(err)
	}

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
