package flags

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MultiString []string

func (m *MultiString) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MultiString) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func (m *MultiString) Contains(value string) bool {

	for _, test := range *m {

		if test == value {
			return true
		}
	}

	return false
}

type MultiDSNString []map[string]string

func (m *MultiDSNString) String() string {

	dsn_strings := make([]string, 0)

	for _, dict := range *m {

		pairs := make([]string, 0)

		for k, v := range dict {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}

		dsn_strings = append(dsn_strings, strings.Join(pairs, " "))
	}

	return strings.Join(dsn_strings, ";;")
}

func (m *MultiDSNString) Set(value string) error {

	value = strings.Trim(value, " ")

	// this is largely so that we can define multiple -dsn-string flags in
	// a single environment variable (20180822/thisisaaronland)

	for _, str_dsn := range strings.Split(value, ";;") {

		str_dsn = strings.Trim(str_dsn, " ")
		pairs := strings.Split(str_dsn, " ")

		dict := make(map[string]string)

		for _, str_pair := range pairs {

			str_pair = strings.Trim(str_pair, " ")
			pair := strings.Split(str_pair, "=")

			if len(pair) != 2 {
				return errors.New("Invalid pair")
			}

			k := pair[0]
			v := pair[1]

			dict[k] = v
		}

		*m = append(*m, dict)
	}

	return nil
}

type MultiInt64 []int64

func (m *MultiInt64) String() string {

	str_values := make([]string, len(*m))

	for i, v := range *m {
		str_values[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(str_values, "\n")
}

func (m *MultiInt64) Set(str_value string) error {

	value, err := strconv.ParseInt(str_value, 10, 64)

	if err != nil {
		return err
	}

	*m = append(*m, value)
	return nil
}

func (m *MultiInt64) Contains(value int64) bool {

	for _, test := range *m {

		if test == value {
			return true
		}
	}

	return false
}
