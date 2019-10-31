package uri

import (
	"fmt"
	_ "log"
	"net/url"
	_ "path/filepath"
)

const FileDriverName string = "file"

type FileURIDriver struct {
	Driver
}

type FileURI struct {
	URI
	path string
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

	// path := filepath.Base(u.Path)
	path := u.Path

	f_u := FileURI{
		path: path,
	}

	return &f_u, nil
}

func (u *FileURI) Driver() string {
	return FileDriverName
}

func (u *FileURI) Origin() string {
	return u.path
}

func (u *FileURI) Target(args ...interface{}) string {
	return u.Origin()
}

func (u *FileURI) String() string {
	return fmt.Sprintf("%s://%s", u.Driver(), u.path)
}
