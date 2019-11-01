package uri

import (
	"errors"
	"fmt"
	"github.com/aaronland/go-string/dsn"
	"github.com/aaronland/go-string/random"
	_ "log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

const IdSecretDriverName string = "idsecret"

func init() {
	dr := NewIdSecretURIDriver()
	RegisterDriver(IdSecretDriverName, dr)
}

type IdSecretURIDriver struct {
	Driver
}

func NewIdSecretURIDriver() Driver {

	dr := IdSecretURIDriver{}
	return &dr
}

func (dr *IdSecretURIDriver) NewURI(str_uri string) (URI, error) {

	return NewIdSecretURI(str_uri)
}

type IdSecretURI struct {
	URI
	origin   string
	id       int64
	secret   string
	secret_o string
	label    string
	format   string
}

func NewIdSecretURIFromDSN(dsn_raw string) (URI, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(dsn_raw, "id", "uri")

	if err != nil {
		return nil, err
	}

	origin := dsn_map["id"]
	id := dsn_map["uri"]

	q := url.Values{}
	q.Set("id", id)

	secret, ok := dsn_map["secret"]

	if ok {
		q.Set("secret", secret)
	}

	secret_o, ok := dsn_map["secret_o"]

	if ok {
		q.Set("secret_o", secret_o)
	}

	raw_uri := fmt.Sprintf("%s?%s", origin, q.Encode())
	str_uri := NewIdSecretURIString(raw_uri)

	return NewIdSecretURI(str_uri)
}

func NewIdSecretURI(str_uri string) (URI, error) {

	u, err := url.Parse(str_uri)

	if err != nil {
		return nil, err
	}

	origin := strings.TrimLeft(u.Path, "/")

	if origin == "" {
		return nil, errors.New("Invalid path")
	}

	q := u.Query()

	str_id := q.Get("id")

	if str_id == "" {
		return nil, errors.New("Missing id")
	}

	id, err := strconv.ParseInt(str_id, 10, 64)

	if err != nil {
		return nil, err
	}

	secret := q.Get("secret")
	secret_o := q.Get("secret_o")
	label := q.Get("label")
	format := q.Get("format")

	rnd_opts := random.DefaultOptions()
	rnd_opts.AlphaNumeric = true

	if secret == "" {

		s, err := random.String(rnd_opts)

		if err != nil {
			return nil, err
		}

		secret = s
	}

	if secret_o == "" {

		s, err := random.String(rnd_opts)

		if err != nil {
			return nil, err
		}

		secret_o = s
	}

	id_u := IdSecretURI{
		origin:   origin,
		id:       id,
		secret:   secret,
		secret_o: secret_o,
		label:    label,
		format:   format,
	}

	return &id_u, nil
}

func (u *IdSecretURI) Driver() string {
	return IdSecretDriverName
}

func (u *IdSecretURI) Target(opts *url.Values) (string, error) {

	str_id := strconv.FormatInt(u.id, 10)

	root := id2Path(u.id)

	secret := u.secret
	format := u.format
	label := u.label

	if opts != nil {

		format = opts.Get("format")
		label = opts.Get("label")
		original := opts.Get("original")

		if original != "" {
			secret = u.secret_o
		}
	}

	if format == "" {
		return "", errors.New("Missing format parameter")
	}

	if label == "" {
		return "", errors.New("Missing label parameter")
	}

	fname := fmt.Sprintf("%s_%s_%s.%s", str_id, secret, label, format)
	uri := filepath.Join(root, fname)

	return uri, nil
}

func (u *IdSecretURI) Origin() string {
	return u.origin
}

func (u *IdSecretURI) String() string {

	q := url.Values{}
	q.Set("id", strconv.FormatInt(u.id, 10))
	q.Set("secret", u.secret)
	q.Set("secret_o", u.secret_o)

	raw_uri := fmt.Sprintf("%s?%s", u.origin, q.Encode())
	return NewIdSecretURIString(raw_uri)
}

func id2Path(id int64) string {

	parts := []string{""}
	input := strconv.FormatInt(id, 10)

	for len(input) > 3 {

		chunk := input[0:3]
		input = input[3:]
		parts = append(parts, chunk)
	}

	if len(input) > 0 {
		parts = append(parts, input)
	}

	return filepath.Join(parts...)
}

func NewIdSecretURIString(raw_uri string) string {
	return fmt.Sprintf("%s:///%s", IdSecretDriverName, raw_uri)
}
