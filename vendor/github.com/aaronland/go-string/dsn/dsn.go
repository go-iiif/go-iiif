// package dsn provides methods for working with data source name (DSN) strings.
package dsn

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// DSN represents a data source name (DSN) style string.
type DSN map[string]string

// Keys returns the list of key names contained in 'dsn'.
func (dsn DSN) Keys() []string {

	keys := make([]string, 0)

	for k, _ := range dsn {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

// String returns a string representation of 'dsn'.
func (dsn DSN) String() string {

	pairs := make([]string, 0)

	for _, k := range dsn.Keys() {

		pair := fmt.Sprintf("%s=%s", k, dsn[k])
		pairs = append(pairs, pair)
	}

	return strings.Join(pairs, " ")
}

// StringToDSN parses 'str_dsn` in to a `DSN` instance.
func StringToDSN(str_dsn string) (DSN, error) {

	str_dsn = strings.Trim(str_dsn, " ")
	dsn := make(map[string]string)

	// see this? it's not particularly smart (translation: not smart at all)
	// about properties with spaces in them - this should be fixed...
	// for example: foo=bar?firstname lastname baz=wubwubwub
	// (20190712/thisisaaronland)

	for _, str_pair := range strings.Split(str_dsn, " ") {

		pair := strings.SplitN(str_pair, "=", 2)

		if len(pair) != 2 {
			return nil, errors.New("Invalid DSN string")
		}

		k := pair[0]
		v := pair[1]

		_, ok := dsn[k]

		if ok {
			msg := fmt.Sprintf("'%s' key already set", k)
			return nil, errors.New(msg)
		}

		dsn[k] = v
	}

	return dsn, nil
}

// StringToDSNWithKeys parse 'str_dsn' in to a `DSN` instance ensuring that the instance contains
// keys matching 'keys'.
func StringToDSNWithKeys(str_dsn string, keys ...string) (DSN, error) {

	dsn, err := StringToDSN(str_dsn)

	if err != nil {
		return nil, err
	}

	for _, k := range keys {

		_, ok := dsn[k]

		if !ok {
			msg := fmt.Sprintf("DSN is missing '%s=' key", k)
			return nil, errors.New(msg)
		}
	}

	return dsn, nil
}
