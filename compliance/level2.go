package compliance

// http://iiif.io/api/image/2.1/compliance/

import (
	"encoding/json"
	"errors"
	"fmt"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"regexp"
)

var level2_spec = `{
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
	     		"bitonal": { "syntax": "bitonal", "required": true, "supported": true, "match": "^bitonal$" }
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
}`

/*
    "http": [
    	    { "feature": "base URI redirects", "name": "baseUriRedirect", "required": true, "supported": 1 },
	    { "feature": "CORS", "name": "cors", "required": true, "supported": 1 },
	    { "feature": "json-ld media type", "name": "jsonldMediaType", "required": true, "supported": 1 },
	    { "feature": "profile link header", "name": "profileLinkHeader", "required": false, "supported": 0 },
	    { "feature": "canonical link header", "name": "canonicalLinkHeader", "required": false, "supported": 0 }
    ]
*/

type Level2ComplianceSpec struct {
	Image ImageCompliance `json:"image"`
}

type Level2Compliance struct {
	Compliance
	spec Level2ComplianceSpec
}

func NewLevel2Compliance(config *iiifconfig.Config) (*Level2Compliance, error) {

	spec := Level2ComplianceSpec{}
	err := json.Unmarshal([]byte(level2_spec), &spec)

	if err != nil {
		return nil, err
	}

	compliance := Level2Compliance{
		spec: spec,
	}

	return &compliance, nil
}

func (c *Level2Compliance) IsValidImageRegion(region string) (bool, error) {

	return c.isvalid(c.spec.Image.Region, region)
}

func (c *Level2Compliance) IsValidImageSize(size string) (bool, error) {

	return c.isvalid(c.spec.Image.Size, size)
}

func (c *Level2Compliance) IsValidImageRotation(rotation string) (bool, error) {

	return c.isvalid(c.spec.Image.Rotation, rotation)
}

func (c *Level2Compliance) IsValidImageQuality(quality string) (bool, error) {

	return c.isvalid(c.spec.Image.Quality, quality)
}

func (c *Level2Compliance) IsValidImageFormat(format string) (bool, error) {

	return c.isvalid(c.spec.Image.Format, format)
}

func (c *Level2Compliance) Formats() []string {

	return c.properties(c.spec.Image.Format)
}

func (c *Level2Compliance) Qualities() []string {

	return c.properties(c.spec.Image.Quality)
}

func (c *Level2Compliance) isvalid(sect map[string]ComplianceDetails, property string) (bool, error) {

	ok := false

	for _, details := range sect {

		if !details.Supported {
			continue
		}

		re, err := regexp.Compile(details.Match)

		if err != nil {
			return false, err
		}

		if re.MatchString(property) {
			ok = true
			break
		}
	}

	if !ok {
		message := fmt.Sprintf("Invalid IIIF 2.1 ...")
		return false, errors.New(message)
	}

	return true, nil
}

func (c *Level2Compliance) properties(sect map[string]ComplianceDetails) []string {

	properties := make([]string, 0)

	for name, details := range sect {

		if !details.Supported {
			continue
		}

		properties = append(properties, name)
	}

	return properties
}

func (c *Level2Compliance) Spec() ([]byte, error) {

	return json.Marshal(c.spec)
}
