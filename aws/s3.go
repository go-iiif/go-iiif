package aws

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-aws/s3"
	"io"
	"io/ioutil"
	"strings"
)

func S3ConfigWrapper(cfg *s3.S3Config) *s3.S3Config {

	if strings.HasPrefix(cfg.Credentials, "shared:") {
		cfg.Credentials = strings.TrimLeft(cfg.Credentials, "shared:")
	}

	return cfg
}

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

func S3SetWrapper(conn *s3.S3Connection, key string, body []byte) error {

	key = fmt.Sprintf("%s#acl=public-read", key)

	r := bytes.NewReader(body)
	fh := ioutil.NopCloser(r)

	return conn.Put(key, fh)
}
