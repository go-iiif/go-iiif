package uri

import (
	"github.com/aaronland/go-string/dsn"
	"github.com/aaronland/go-string/random"
)

type IdSecretURI struct {
	URI
	dsn_map dsn.DSN
}

func NewIdSecretURI(raw string) (URI, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(raw, "id", "uri")

	if err != nil {
		return nil, err
	}

	opts := random.DefaultOptions()
	opts.AlphaNumeric = true

	_, ok := dsn_map["secret"]

	if !ok {

		s, err := random.String(opts)

		if err != nil {
			return nil, err
		}

		dsn_map["secret"] = s
	}

	_, ok = dsn_map["secret_o"]

	if !ok {

		s, err := random.String(opts)

		if err != nil {
			return nil, err
		}

		dsn_map["secret_o"] = s
	}

	u := IdSecretURI{
		dsn_map: dsn_map,
	}

	return &u, nil
}

func (u *IdSecretURI) URL() string {
	url, _ := u.dsn_map["uri"]
	return url
}

func (u *IdSecretURI) String() string {
	return u.dsn_map.String()
}
