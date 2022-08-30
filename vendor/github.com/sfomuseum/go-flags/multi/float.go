package multi

import (
	"strconv"
	"strings"
)

type MultiFloat64 []float64

func (m *MultiFloat64) String() string {

	str_values := make([]string, len(*m))

	for i, v := range *m {
		str_values[i] = strconv.FormatFloat(v, 'f', 10, 64)
	}

	return strings.Join(str_values, "\n")
}

func (m *MultiFloat64) Set(str_value string) error {

	value, err := strconv.ParseFloat(str_value, 64)

	if err != nil {
		return err
	}

	*m = append(*m, value)
	return nil
}

func (m *MultiFloat64) Get() interface{} {
	return *m
}

func (m *MultiFloat64) Contains(value float64) bool {

	for _, test := range *m {

		if test == value {
			return true
		}
	}

	return false
}
