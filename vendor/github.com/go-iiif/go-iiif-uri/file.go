package uri

import (
	"fmt"
	"net/url"
)

const FileDriverName string = "file"

type FileURIDriver struct {
	Driver
}

type FileURI struct {
	URI
	origin string
	target string
}

func init() {
	dr := NewFileURIDriver()
	RegisterDriver(FileDriverName, dr)
}

func NewFileURIDriver() Driver {

	dr := FileURIDriver{}
	return &dr
}

func (dr *FileURIDriver) NewURI(str_uri string) (URI, error) {
	return NewFileURI(str_uri)
}

func NewFileURI(str_uri string) (URI, error) {

	u, err := url.Parse(str_uri)

	if err != nil {
		return nil, err
	}

	origin := u.Path

	q := u.Query()

	target := q.Get("target")

	if target == "" {
		target = origin
	}

	f_u := FileURI{
		origin: origin,
		target: target,
	}

	return &f_u, nil
}

func (u *FileURI) Driver() string {
	return FileDriverName
}

func (u *FileURI) Origin() string {
	return u.origin
}

func (u *FileURI) Target(opts *url.Values) (string, error) {
	return u.target, nil
}

func (u *FileURI) String() string {

	str_uri := fmt.Sprintf("%s://%s", u.Driver(), u.origin)

	if u.target != "" && u.target != u.origin {
		q := url.Values{}
		q.Set("target", u.target)
		str_uri = fmt.Sprintf("%s?%s", str_uri, q.Encode())
	}

	return str_uri
}
