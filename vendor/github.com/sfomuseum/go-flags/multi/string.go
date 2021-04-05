package multi

import (
	"strings"
)

type MultiString []string

func (m *MultiString) String() string {
	return strings.Join(*m, "\n")
}

func (m *MultiString) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func (m *MultiString) Get() interface{} {
	return *m
}

func (m *MultiString) Contains(value string) bool {

	for _, test := range *m {

		if test == value {
			return true
		}
	}

	return false
}
