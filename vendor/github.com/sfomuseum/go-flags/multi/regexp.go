package multi

import (
	"fmt"
	"regexp"
	"strings"
)

type MultiRegexp []*regexp.Regexp

func (i *MultiRegexp) String() string {

	patterns := make([]string, 0)

	for _, re := range *i {
		patterns = append(patterns, fmt.Sprintf("%v", re))
	}

	return strings.Join(patterns, "\n")
}

func (i *MultiRegexp) Set(value string) error {

	re, err := regexp.Compile(value)

	if err != nil {
		return err
	}

	*i = append(*i, re)
	return nil
}

func (i *MultiRegexp) Get() interface{} {
	return *i
}
