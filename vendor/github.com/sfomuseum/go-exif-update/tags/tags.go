package tags

import (
	_ "embed"
	"github.com/dsoprea/go-exif/v3/common"
	"gopkg.in/yaml.v2"
	_ "log"
	"sync"
)

// tags_data.yaml has been cloned from:
// https://raw.githubusercontent.com/dsoprea/go-exif/de2141190595193aa097a2bf3205ba0cf76dc14b/tags_data.go

//go:embed tags_data.yaml
var tags_data []byte

// sync.Once instance for loading YAML-encoded tag data.
var tags_init sync.Once

// tags_suppported is a private list of tag names supported by this package. This is used in conjunction with the tags_init sync.Once instance.
var tags_supported []string

// A list of exifcommon.TagTypePrimitive instances for tags that are not supported yet by this package.
var UnsupportedTypes []exifcommon.TagTypePrimitive

// A list of string for EXIF tag types that are not supported yet by this package.
var UnsupportedTypesString []string

func init() {

	// These are not supported yet because I am not sure what if any value-wrangling
	// we need to do to ensure they are recorded as valid EXIF tags.
	// (20210409/thisisaaronland)

	UnsupportedTypes = []exifcommon.TagTypePrimitive{
		// exifcommon.TypeRational,
		// exifcommon.TypeSignedRational,
		exifcommon.TypeShort,
		exifcommon.TypeLong,
		exifcommon.TypeSignedLong,
		exifcommon.TypeUndefined,
	}

	unsupported_str := make([]string, len(UnsupportedTypes))

	for idx, t := range UnsupportedTypes {
		unsupported_str[idx] = t.String()
	}

	UnsupportedTypesString = unsupported_str
}

// encodedTag is a struct for holding YAML-encoded EXIF tags. This is cribbed from here because it's a private type:
// https://github.com/dsoprea/go-exif/blob/de2141190595193aa097a2bf3205ba0cf76dc14b/tags.go#L189
type encodedTag struct {
	// id is signed, here, because YAML doesn't have enough information to
	// support unsigned.
	Id       int    `yaml:"id"`
	Name     string `yaml:"name"`
	TypeName string `yaml:"type_name"`
}

// Determine whether a string tag name is included in the list of supported tags in this package.
func IsSupported(t string) (bool, error) {

	supported, err := SupportedTags()

	if err != nil {
		return false, err
	}

	for _, this_t := range supported {

		if this_t == t {
			return true, nil
		}
	}

	return false, nil
}

// Return a list of tag names that are supported by this package.
func SupportedTags() ([]string, error) {

	var tags_err error

	tags_func := func() {

		tags_supported = make([]string, 0)

		encodedIfds := make(map[string][]encodedTag)

		err := yaml.Unmarshal(tags_data, encodedIfds)

		if err != nil {
			tags_err = err
			return
		}

		for _, ifdtags := range encodedIfds {

			for _, t := range ifdtags {

				include := true

				for _, ts := range UnsupportedTypesString {

					if t.TypeName == ts {
						include = false
						break
					}
				}

				if include {
					tags_supported = append(tags_supported, t.Name)
				}
			}
		}
	}

	tags_init.Do(tags_func)

	if tags_err != nil {
		return nil, tags_err
	}

	return tags_supported, nil
}
