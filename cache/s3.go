package cache

import (
	"fmt"

	_ "github.com/aaronland/gocloud-blob/s3"
	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
)

// NewS3CacheURIFromConfig returns a URI derived from 'cfg'.
func NewS3CacheURIFromConfig(cfg iiifconfig.CacheConfig) (string, error) {

	bucket := cfg.Path
	prefix := cfg.Prefix
	region := cfg.Region
	creds := cfg.Credentials

	uri := fmt.Sprintf("s3blob://%s?region=%s&credentials=%s&prefix=%s", bucket, region, creds, prefix)
	return uri, nil
}

// NewS3Cache returns a NewBlobCacheFromURI with a constructed blob uri.
func NewS3Cache(cfg iiifconfig.CacheConfig) (Cache, error) {

	uri := cfg.URI

	if uri == "" {
		v, err := NewS3CacheURIFromConfig(cfg)

		if err != nil {
			return nil, err
		}

		uri = v
	}

	return NewBlobCacheFromURI(uri)
}

func NewS3CacheFromURI(uri string) (Cache, error) {
	return NewBlobCacheFromURI(uri)
}
