package compliance

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
)

// http://iiif.io/api/image/2.1/
// http://iiif.io/api/image/2.1/compliance/

//go:embed level2.json
var level2_spec []byte

type Level2Compliance struct {
	Compliance
	spec *ComplianceSpec
}

func NewLevel2Compliance(config *iiifconfig.Config) (*Level2Compliance, error) {

	spec, err := NewLevel2ComplianceSpecWithConfig(config)

	if err != nil {
		return nil, err
	}

	compliance := Level2Compliance{
		spec: spec,
	}

	return &compliance, nil
}

func NewLevel2ComplianceSpec() (*ComplianceSpec, error) {

	spec := ComplianceSpec{}
	err := json.Unmarshal(level2_spec, &spec)

	if err != nil {
		return nil, err
	}

	return &spec, nil
}

func NewLevel2ComplianceSpecWithConfig(config *iiifconfig.Config) (*ComplianceSpec, error) {

	spec, err := NewLevel2ComplianceSpec()

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

	return spec, err
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

func (c *Level2Compliance) Supports() []string {

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

		re, err := regexp.Compile(details.Match)

		if err != nil {
			return false, err
		}

		if !re.MatchString(value) {
			continue
		}

		if !details.Supported {
			message := fmt.Sprintf("Unsupported IIIF 2.1 feature (%s) '%s'", name, value)
			return false, errors.New(message)
		}

		ok = true
		break

	}

	if !ok {
		message := fmt.Sprintf("Invalid IIIF 2.1 feature property (%s) : '%s'", property, value)
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

func (c *Level2Compliance) Spec() *ComplianceSpec {

	return c.spec
}

func (c *Level2Compliance) DefaultQuality() (string, error) {

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
