package compliance

// http://iiif.io/api/image/2.1/
// http://iiif.io/api/image/2.1/compliance/

import (
	"encoding/json"
	"errors"
	"fmt"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	_ "log"
	"regexp"
)

var level2_spec = `{
    "image": {
    	     "region": {
	     	       "full":         { "syntax": "full",        "required": true, "supported": true, "match": "^full$" },
		       "regionByPx":   { "syntax": "x,y,w,h",     "required": true, "supported": true, "match": "^\\d+\\,\\d+\\,\\d+\\,\\d+$" },
		       "regionByPct":  { "syntax": "pct:x,y,w,h", "required": true, "supported": true, "match": "^pct\\:\\d+\\,\\d+\\,\\d+\\,\\d+$" },
		       "regionSquare": { "syntax": "square",      "required": false, "supported": false, "match": "^square$" }
	     },
	     "size": {
	     		"full":              { "syntax": "full",  "required": true, "supported": true, "match": "^full$" },
	     		"max":               { "syntax": "max",   "required": false, "supported": true, "match": "^max$" },
	     		"sizeByW":           { "syntax": "w,",    "required": true, "supported": true, "match": "^\\d+\\,$" },			
	     		"sizeByH":           { "syntax": ",h",    "required": true, "supported": true, "match": "^\\,\\d+$" },
	     		"sizeByPct":         { "syntax": "pct:n", "required": true, "supported": true, "match": "^pct\\:\\d+(\\.\\d+)?$" },			
	     		"sizeByConfinedWh":  { "syntax": "!w,h",  "required": true, "supported": true, "match": "^\\!\\d+\\,\\d+$" },
	     		"sizeByDistortedWh": { "syntax": "w,h",   "required": true, "supported": true, "match": "^\\d+\\,\\d+$" },
	     		"sizeByWh":          { "syntax": "w,h",   "required": true, "supported": true, "match": "^\\d+\\,\\d+$" }
	     },
	     "rotation": {
	     		"none":              { "syntax": "0",          "required": true, "supported": true, "match": "^0$" },
	     		"rotationBy90s":     { "syntax": "90,180,270", "required": true, "supported": true, "match": "^(?:90|180|270)$" },
	     		"rotationArbitrary": { "syntax": "",           "required": false, "supported": true, "match": "^\\d+\\.\\d+$" },			
	     		"mirroring":         { "syntax": "!n",         "required": true, "supported": true, "match": "^\\!\\d+$" }
	     },
	     "quality": {
	     		"default": { "syntax": "default", "required": true, "supported": true, "match": "^default$" },
	     		"color":   { "syntax": "color",   "required": false, "supported": true, "match": "^colou?r$" },
	     		"gray":    { "syntax": "gray",    "required": false, "supported": false, "match": "gr(?:e|a)y$" },			
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
	spec *Level2ComplianceSpec
}

func NewLevel2Compliance(config *iiifconfig.Config) (*Level2Compliance, error) {

	spec, err := NewLevel2ComplianceSpec(config)

	if err != nil {
		return nil, err
	}

	compliance := Level2Compliance{
		spec: spec,
	}

	return &compliance, nil
}

func NewLevel2ComplianceSpec(config *iiifconfig.Config) (*Level2ComplianceSpec, error) {

	spec := Level2ComplianceSpec{}
	err := json.Unmarshal([]byte(level2_spec), &spec)

	if err != nil {
		return nil, err
	}

	feature_block := func(block string) (map[string]ComplianceDetails, error) {

		var possible map[string]ComplianceDetails

		if block == "region" {
			possible = spec.Image.Region
		} else if block == "size" {
			possible = spec.Image.Size
		} else if block == "rotation" {
			possible = spec.Image.Rotation
		} else if block == "quality" {
			possible = spec.Image.Quality
		} else if block == "format" {
			possible = spec.Image.Format
		} else {
			message := fmt.Sprintf("Unknown block %s", block)
			return nil, errors.New(message)
		}

		return possible, nil
	}

	toggle_features := func(stuff iiifconfig.FeaturesToggle, toggle bool) error {

		for block, features := range stuff {

			possible, err := feature_block(block)

			if err != nil {
				return err
			}

			for _, f := range features {

				details, ok := possible[f]

				if !ok {
					message := fmt.Sprintf("Undefined feature %s for block (%s)", f, block)
					return errors.New(message)
				}

				details.Supported = toggle
				possible[f] = details
			}
		}

		return nil
	}

	append_features := func(stuff iiifconfig.FeaturesAppend) error {

		for block, features := range stuff {

			possible, err := feature_block(block)

			if err != nil {
				return err
			}

			for name, details := range features {

				possible[name] = ComplianceDetails{
					Syntax:    details.Syntax,
					Required:  details.Required,
					Supported: details.Supported,
					Match:     details.Match,
				}
			}
		}

		return nil
	}

	err = append_features(config.Features.Append)

	if err != nil {
		return nil, err
	}

	err = toggle_features(config.Features.Enable, true)

	if err != nil {
		return nil, err
	}

	err = toggle_features(config.Features.Disable, false)

	if err != nil {
		return nil, err
	}

	return &spec, err
}

func (c *Level2Compliance) IsValidImageRegion(region string) (bool, error) {

	return c.isvalid("region", region)
}

func (c *Level2Compliance) IsValidImageSize(size string) (bool, error) {

	return c.isvalid("size", size)
}

func (c *Level2Compliance) IsValidImageRotation(rotation string) (bool, error) {

	return c.isvalid("rotation", rotation)
}

func (c *Level2Compliance) IsValidImageQuality(quality string) (bool, error) {

	return c.isvalid("quality", quality)
}

func (c *Level2Compliance) IsValidImageFormat(format string) (bool, error) {

	return c.isvalid("format", format)
}

func (c *Level2Compliance) Formats() []string {

	return c.properties(c.spec.Image.Format)
}

func (c *Level2Compliance) Qualities() []string {

	return c.properties(c.spec.Image.Quality)
}

func (c *Level2Compliance) isvalid(property string, value string) (bool, error) {

	var sect map[string]ComplianceDetails

	if property == "region" {
		sect = c.spec.Image.Region
	} else if property == "size" {
		sect = c.spec.Image.Size
	} else if property == "rotation" {
		sect = c.spec.Image.Rotation
	} else if property == "quality" {
		sect = c.spec.Image.Quality
	} else if property == "format" {
		sect = c.spec.Image.Format
	} else {
		message := fmt.Sprintf("Unknown property %s", property)
		return false, errors.New(message)
	}

	ok := false

	for name, details := range sect {

		// log.Printf("%s %t (%s = %s)", name, details.Supported, property, value)

		re, err := regexp.Compile(details.Match)

		if err != nil {
			return false, err
		}

		if !re.MatchString(value) {
			continue
		}

		if !details.Supported {
			message := fmt.Sprintf("Unsupported IIIF 2.1 feature (%s) %s", name, value)
			return false, errors.New(message)
		}

		// log.Printf("%s %s MATCH %s", name, property, value)
		ok = true
		break

	}

	if !ok {
		message := fmt.Sprintf("Invalid IIIF 2.1 feature property %s %s", property, value)
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
