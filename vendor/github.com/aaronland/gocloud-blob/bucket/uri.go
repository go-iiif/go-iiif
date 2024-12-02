package bucket

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// Parse a URI in to its `gocloud.dev/blob.Bucket` URI and key values.
func ParseURI(uri string) (string, string, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return "", "", fmt.Errorf("Failed to parse source URI, %w", err)
	}

	var bucket_uri string
	var bucket_key string

	if u.Scheme == "" {

		path := u.Path
		abs_path, err := filepath.Abs(path)

		if err != nil {
			return "", "", fmt.Errorf("Failed to derive absolute path, %w", err)
		}

		root := filepath.Dir(abs_path)
		base := filepath.Base(abs_path)

		root = strings.Trim(root, "/")
		root = fmt.Sprintf("%s/", root)

		q := u.Query()
		q.Set("prefix", root)

		u.Scheme = "file"
		u.Host = u.Host
		u.Path = "/"

		u.RawQuery = q.Encode()

		bucket_uri = u.String()
		bucket_key = base

	} else {

		path := u.Path
		root := filepath.Dir(path)
		base := filepath.Base(path)

		root = strings.TrimRight(root, "/")
		root = fmt.Sprintf("%s/", root)

		u.Path = root

		bucket_uri = u.String()
		bucket_key = base
	}

	return bucket_uri, bucket_key, nil
}
