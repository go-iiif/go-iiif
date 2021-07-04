package uri

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronland/go-string/random"
	_ "log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

const IDSECRET_SCHEME string = "idsecret"

type IdSecretURI struct {
	URI
	origin   string
	id       string
	secret   string
	secret_o string
	label    string
	format   string
	prefix   string
}

func init() {
	ctx := context.Background()
	RegisterURI(ctx, IDSECRET_SCHEME, NewIdSecretURI)
}

func NewIdSecretURI(ctx context.Context, str_uri string) (URI, error) {

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

	if q.Get("ensure-int") != "" {

		_, err := strconv.ParseUint(str_id, 10, 64)

		if err != nil {
			return nil, err
		}
	}

	secret := q.Get("secret")
	secret_o := q.Get("secret_o")
	label := q.Get("label")
	format := q.Get("format")
	prefix := q.Get("prefix")

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
		id:       str_id,
		secret:   secret,
		secret_o: secret_o,
		label:    label,
		format:   format,
		prefix:   prefix,
	}

	return &id_u, nil
}

func (u *IdSecretURI) Target(opts *url.Values) (string, error) {

	str_id := u.id // strconv.FormatUint(u.id, 10)

	prefix := Id2Path(u.id)

	if u.prefix != "" {
		prefix = u.prefix
	}

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
	uri := filepath.Join(prefix, fname)

	return uri, nil
}

func (u *IdSecretURI) Origin() string {
	return u.origin
}

func (u *IdSecretURI) String() string {

	q := url.Values{}
	q.Set("id", u.id) // strconv.FormatUint(u.id, 10))
	q.Set("secret", u.secret)
	q.Set("secret_o", u.secret_o)

	raw_uri := fmt.Sprintf("%s?%s", u.origin, q.Encode())

	return fmt.Sprintf("%s:///%s", u.Scheme(), raw_uri)
}

func (u *IdSecretURI) Scheme() string {
	return IDSECRET_SCHEME
}

func Id2Path(id string) string {

	parts := []string{""}
	input := id // strconv.FormatUint(id, 10)

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
