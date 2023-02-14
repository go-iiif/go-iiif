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

type MultiCSVString []string

func (m *MultiCSVString) String() string {
	return strings.Join(*m, "\n")
}

func (m *MultiCSVString) Set(value string) error {

	for _, v := range strings.Split(value, ",") {
		*m = append(*m, v)
	}

	return nil
}

func (m *MultiCSVString) Get() interface{} {
	return *m
}

func (m *MultiCSVString) Contains(value string) bool {

	for _, test := range *m {

		if test == value {
			return true
		}
	}

	return false
}
