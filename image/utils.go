package image

import (
	"net/url"
	"strings"
)

/*

things I am unsure of (20160911/thisisaaronland)
1. whether this should call sanitize.SanitizeString
2. whether this should even exist in its own package

*/

func ScrubIdentifier(identifier string) (string, error) {

	clean, err := url.QueryUnescape(identifier)

	if err != nil {
		return "", err
	}

	clean = strings.Replace(clean, "../", "", -1)
	return clean, nil
}
