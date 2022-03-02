package compliance

import (
	"encoding/json"
	"errors"
	"fmt"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	_ "log"
	"regexp"
)

var level0_spec = `{
    "image": {
 	     "region": {},
	     "size": {},
	     "rotation": {},
	     "quality": {
	     		"default": { "syntax": "default", "required": true, "supported": true, "match": "^default$", "default": true }
             },
	     "format": {
	     	       "jpg": { "syntax": "jpg",  "required": true, "supported": true, "match": "^jpe?g$" }
	     }	     
    },
    "http": {
            "baseUriRedirect":     { "name": "base URI redirects",    "required": true,  "supported": true },
	    "cors":                { "name": "CORS",                  "required": true,  "supported": true },
	    "jsonldMediaType":     { "name": "json-ld media type",    "required": true,  "supported": true },
	    "profileLinkHeader":   { "name": "profile link header",   "required": false, "supported": false },
	    "canonicalLinkHeader": { "name": "canonical link header", "required": false, "supported": false }
    }
}`

type Level0Compliance struct {
	Compliance
	spec *ComplianceSpec
}

func NewLevel0Compliance(config *iiifconfig.Config) (*Level0Compliance, error) {

	spec, err := NewLevel0ComplianceSpecWithConfig(config)

	if err != nil {
		return nil, err
	}

	compliance := Level0Compliance{
		spec: spec,
	}

	return &compliance, nil
}

func NewLevel0ComplianceSpec() (*ComplianceSpec, error) {

	spec := ComplianceSpec{}
	err := json.Unmarshal([]byte(level0_spec), &spec)

	if err != nil {
		return nil, err
	}

	return &spec, nil
}

func NewLevel0ComplianceSpecWithConfig(config *iiifconfig.Config) (*ComplianceSpec, error) {

	spec, err := NewLevel0ComplianceSpec()

	if err != nil {
		return nil, err
	}

	return spec, nil
}

func (c *Level0Compliance) IsValidImageRegion(region string) (bool, error) {

	return c.isvalid("region", region)
}

func (c *Level0Compliance) IsValidImageSize(size string) (bool, error) {

	return c.isvalid("size", size)
}

func (c *Level0Compliance) IsValidImageRotation(rotation string) (bool, error) {

	return c.isvalid("rotation", rotation)
}

func (c *Level0Compliance) IsValidImageQuality(quality string) (bool, error) {

	return c.isvalid("quality", quality)
}

func (c *Level0Compliance) IsValidImageFormat(format string) (bool, error) {

	return c.isvalid("format", format)
}

func (c *Level0Compliance) Formats() []string {

	return c.properties(c.spec.Image.Format)
}

func (c *Level0Compliance) Qualities() []string {

	return c.properties(c.spec.Image.Quality)
}

func (c *Level0Compliance) Supports() []string {

	supports := make([]string, 0)

	for _, s := range c.properties(c.spec.Image.Region) {
		supports = append(supports, s)
	}

	for _, s := range c.properties(c.spec.Image.Size) {
		supports = append(supports, s)
	}

	for _, s := range c.properties(c.spec.Image.Rotation) {
		supports = append(supports, s)
	}

	for name, details := range c.spec.HTTP {

		if !details.Supported {
			continue
		}

		supports = append(supports, name)
	}

	return supports
}

func (c *Level0Compliance) isvalid(property string, value string) (bool, error) {

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

func (c *Level0Compliance) properties(sect map[string]ComplianceDetails) []string {

	properties := make([]string, 0)

	for name, details := range sect {

		if !details.Supported {
			continue
		}

		properties = append(properties, name)
	}

	return properties
}

func (c *Level0Compliance) Spec() *ComplianceSpec {

	return c.spec
}

func (c *Level0Compliance) DefaultQuality() (string, error) {

	quality := ""

	for q, details := range c.spec.Image.Quality {

		if details.Supported && details.Default {
			quality = q
			break
		}

	}

	if quality == "" {
		return "", errors.New("Unable to determine default quality")
	}

	return quality, nil
}
