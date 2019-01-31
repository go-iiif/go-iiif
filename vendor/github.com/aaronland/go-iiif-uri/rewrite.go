package uri

import (
	"github.com/aaronland/go-string/dsn"
	"github.com/aaronland/go-string/random"
)

type RewriteURI struct {
	URI
	dsn_map dsn.DSN
}

func NewRewriteURI(raw string) (URI, error) {

	dsn_map, err := dsn.StringToDSN(raw)

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

	u := RewriteURI{
		dsn_map: dsn_map,
	}

	return &u, nil
}

func (u *RewriteURI) URL() string {
	url, _ := u.dsn_map["uri"]
	return url
}

func (u *RewriteURI) String() string {
	return u.dsn_map.String()
}
