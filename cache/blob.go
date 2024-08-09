package cache

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aaronland/gocloud-blob/bucket"
	aa_s3 "github.com/aaronland/gocloud-blob/s3"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	"gocloud.dev/blob"
)

// A BlobCache is a Cache that includes attributes specific to a blob such as the bucket_uri and acl configuration.
type BlobCache struct {
	Cache
	bucket_uri string
	bucket     *blob.Bucket
	scheme     string
	acl        string
}

func init() {
	ctx := context.Background()
	err := RegisterBlobCacheSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterBloblCacheSchemes will explicitly register all the schemes associated with the
func RegisterBlobCacheSchemes(ctx context.Context) error {

	register_mu.Lock()
	defer register_mu.Unlock()

	to_register := make([]string, 0)

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {
		to_register = append(to_register, scheme)
	}

	for _, scheme := range to_register {

		_, exists := register_map[scheme]

		if exists {
			continue
		}

		err := RegisterCache(ctx, scheme, NewBlobCacheFromURI)

		if err != nil {
			return fmt.Errorf("Failed to register blob cache for '%s', %w", scheme, err)
		}

		register_map[scheme] = true
	}

	return nil
}

// NewBlobCacheURIFromConfig returns a valid cache.Cache URI derived from 'config'.
func NewBlobCacheURIFromConfig(config iiifconfig.CacheConfig) (string, error) {
	return config.Path, nil
}

// NewBlobCache returns a NewBlobCacheFromURI.
func NewBlobCache(cfg iiifconfig.CacheConfig) (Cache, error) {

	uri := cfg.URI

	if uri == "" {
		v, err := NewBlobCacheURIFromConfig(cfg)

		if err != nil {
			return nil, err
		}

		uri = v
	}

	return NewBlobCacheFromURI(uri)
}

// NewBlobCacheFromURI returns a BlobCache using the GoCloud package.
func NewBlobCacheFromURI(uri string) (Cache, error) {

	ctx := context.Background()

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open bucket for %s, %w", uri, err)
	}

	// something something something permissions and ACLs in Go Cloud
	// basically we need to trap a `acl=VALUE` query parameter in order
	// to set permission - as of this writing we a) only handle S3 and
	// b) we assign ACL values at the bucket level and not for specific
	// blobs, at least not yet (20191113/thisisaaronland)
	//
	// https://github.com/google/go-cloud/issues/1108
	// https://godoc.org/gocloud.dev/blob

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	scheme := u.Scheme
	acl := q.Get("acl")

	bc := &BlobCache{
		bucket_uri: uri,
		bucket:     b,
		scheme:     scheme,
		acl:        acl,
	}

	return bc, nil
}

func (bc *BlobCache) String() string {
	return bc.bucket_uri
}

// Exists returns a bool set to true if the configured bucket exists.
func (bc *BlobCache) Exists(uri string) bool {

	ctx := context.Background()

	exists, err := bc.bucket.Exists(ctx, uri)

	if err != nil {
		return false
	}

	return exists
}

// Get reads data from a BlobCache location.
func (bc *BlobCache) Get(uri string) ([]byte, error) {

	ctx := context.Background()

	fh, err := bc.bucket.NewReader(ctx, uri, nil)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return io.ReadAll(fh)
}

// Set writes data to a BlobCache location.
func (bc *BlobCache) Set(uri string, body []byte) error {

	ctx := context.Background()

	var wr_opts *blob.WriterOptions

	// see notes above in NewBlobCacheFromURI

	if strings.HasPrefix(bc.scheme, "s3") && bc.acl != "" {

		acl, err := aa_s3.StringACLToObjectCannedACL(bc.acl)

		if err != nil {
			return fmt.Errorf("Failed to derive ACL object, %w", err)
		}

		before := func(asFunc func(interface{}) bool) error {

			req := &aws_s3.PutObjectInput{}
			ok := asFunc(&req)

			if !ok {
				return errors.New("invalid s3 type")
			}

			req.ACL = acl
			return nil
		}

		wr_opts = &blob.WriterOptions{
			BeforeWrite: before,
		}
	}

	fh, err := bc.bucket.NewWriter(ctx, uri, wr_opts)

	if err != nil {
		return err
	}

	_, err = fh.Write(body)

	if err != nil {
		fh.Close()
		return err
	}

	err = fh.Close()

	if err != nil {
		return err
	}

	return nil
}

// Unset deletes data from a BlobCache location.
func (bc *BlobCache) Unset(uri string) error {
	ctx := context.Background()
	return bc.bucket.Delete(ctx, uri)
}
