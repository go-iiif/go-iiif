package funcs

import (
	"log"
	"net/url"
)

// JoinPath returns the value of 'prefix' prepended to 'path'.
func JoinPath(prefix string, path string) string {

	if prefix == "" {
		return path
	}

	p, err := url.JoinPath(prefix, path)

	if err != nil {
		log.Printf("Failed to join path, %v", err)
		return "#"
	}

	return p
}
