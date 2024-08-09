package bucket

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	gc_blob "gocloud.dev/blob"
)

// OpenBucket is a local helper function to open a gocloud.dev/blob Bucket URI and ensuring
// that files will not be written with their corresponding metdata (`.attrs`) files.
func OpenBucket(ctx context.Context, bucket_uri string) (*gc_blob.Bucket, error) {

	u, err := url.Parse(bucket_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse bucket URI, %w", err)
	}

	if u.Scheme == "cwd" {

		cwd, err := os.Getwd()

		if err != nil {
			return nil, fmt.Errorf("Failed to derive current working directory, %w", err)
		}

		u.Scheme = "file"
		u.Path = filepath.Join(cwd, u.Path)
	}

	if allowsToSkipMetadata(u) {

		q := u.Query()

		if q.Get("metadata") != "skip" {
			q.Set("metadata", "skip")
			u.RawQuery = q.Encode()
		}
	}

	bucket_uri = u.String()

	return gc_blob.OpenBucket(ctx, bucket_uri)
}

// allowsToSkipMetadata returns a boolean value indicating whether or not the blob opener for
// 'u' supports the '?metadata=skip' parameter
func allowsToSkipMetadata(u *url.URL) bool {

	allowed := false

	allowed_schemes := []string{
		"file",
	}

	scheme := u.Scheme

	for _, s := range allowed_schemes {
		if scheme == s {
			allowed = true
			break
		}
	}

	return allowed
}
