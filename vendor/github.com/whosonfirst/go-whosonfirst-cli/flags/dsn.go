package flags

import (
	"errors"
	"strings"
)

// PLEASE RECONCILE ME WITH MultiDSNString and dsn/DSN...
// (20181105/thisisaaronland)

type DSNString [][]string

func (m *DSNString) String() string {

	dsn_strings := make([]string, 0)

	for _, pair := range *m {

		dsn_strings = append(dsn_strings, strings.Join(pair, "="))
	}

	return strings.Join(dsn_strings, " ")
}

func (m *DSNString) Set(value string) error {

	value = strings.Trim(value, " ")

	pairs := strings.Split(value, " ")

	for _, str_pair := range pairs {

		str_pair = strings.Trim(str_pair, " ")
		pair := strings.Split(str_pair, "=")

		if len(pair) != 2 {
			return errors.New("Invalid pair")
		}

		*m = append(*m, pair)
	}

	return nil
}

func (m *DSNString) Map() map[string]string {

	dsn_map := make(map[string]string)

	for _, pair := range *m {

		dsn_map[pair[0]] = pair[1]
	}

	return dsn_map
}
