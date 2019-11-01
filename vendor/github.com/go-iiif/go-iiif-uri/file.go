package uri

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
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

	origin := strings.TrimLeft(u.Path, "/")

	if origin == "" {
		return nil, errors.New("Invalid path")
	}

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

	raw_uri := fmt.Sprintf("%s", u.origin)

	if u.target != "" && u.target != u.origin {
		q := url.Values{}
		q.Set("target", u.target)
		raw_uri = fmt.Sprintf("%s?%s", raw_uri, q.Encode())
	}

	return NewFileURIString(raw_uri)
}

func NewFileURIString(str_uri string) string {
	return fmt.Sprintf("%s:///%s", FileDriverName, str_uri)
}
