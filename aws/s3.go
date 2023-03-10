package aws

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-aws/s3"
)

// S3ConfigWrapper returns an S3Config without the `shared` prefix
func S3ConfigWrapper(cfg *s3.S3Config) *s3.S3Config {

	if strings.HasPrefix(cfg.Credentials, "shared:") {
		cfg.Credentials = strings.TrimLeft(cfg.Credentials, "shared:")
	}

	return cfg
}

// S3GetWrapper fetches an object from AWS S3 and returns bytes or an error
func S3GetWrapper(conn *s3.S3Connection, key string) ([]byte, error) {

	fh, err := conn.Get(key)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	var b bytes.Buffer
	out := bufio.NewWriter(&b)

	_, err = io.Copy(out, fh)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// S3SetWrapper uploads an object to S3 with the `acl=public-read` setting
func S3SetWrapper(conn *s3.S3Connection, key string, body []byte) error {

	key = fmt.Sprintf("%s#acl=public-read", key)

	r := bytes.NewReader(body)
	fh := ioutil.NopCloser(r)

	return conn.Put(key, fh)
}
