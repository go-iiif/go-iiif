package funcs

import (
	"strings"
)

func StringHasPrefix(input string, prefix string) bool {
	return strings.HasPrefix(input, prefix)
}
