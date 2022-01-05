package mimetypes

import (
	"github.com/aaronland/go-mimetypes/lookup"
	"strings"
)

func TypesByExtension(e string) []string {

	if strings.HasPrefix(e, ".") {
		e = strings.Replace(e, ".", "", 1)
	}

	t, ok := lookup.ByEXTENSION[e]

	if ok {
		return t
	}

	return []string{}
}

func ExtensionsByType(t string) []string {

	e, ok := lookup.ByMIMETYPE[t]

	if ok {
		return e
	}

	return []string{}
}
