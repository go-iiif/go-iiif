package uri

import (
	"errors"
	"fmt"
	"net/url"
)

const RewriteDriverName string = "rewrite"

type RewriteURIDriver struct {
	Driver
}

type RewriteURI struct {
	URI
	origin string
	target string
}

func init() {
	dr := NewRewriteURIDriver()
	RegisterDriver(RewriteDriverName, dr)
}

func NewRewriteURIDriver() Driver {

	dr := RewriteURIDriver{}
	return &dr
}

func (dr *RewriteURIDriver) NewURI(str_uri string) (URI, error) {
	return NewRewriteURI(str_uri)
}

func NewRewriteURI(str_uri string) (URI, error) {

	u, err := url.Parse(str_uri)

	if err != nil {
		return nil, err
	}

	origin := u.Path

	q := u.Query()

	target := q.Get("target")

	if target == "" {
		return nil, errors.New("Missing rewrite target")
	}

	if target == origin {
		return nil, errors.New("Invalid rewrite target")
	}

	rw := RewriteURI{
		origin: origin,
		target: target,
	}

	return &rw, nil
}

func (u *RewriteURI) Driver() string {
	return RewriteDriverName
}

func (u *RewriteURI) Origin() string {
	return u.origin
}

func (u *RewriteURI) Target(opts *url.Values) (string, error) {
	return u.target, nil
}

func (u *RewriteURI) String() string {

	str_uri := fmt.Sprintf("%s://%s", u.Driver(), u.origin)

	q := url.Values{}
	q.Set("target", u.target)
	str_uri = fmt.Sprintf("%s?%s", str_uri, q.Encode())

	return str_uri
}
