package cache

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"

	"github.com/aaronland/gocloud-blob/bucket"
	aa_s3 "github.com/aaronland/gocloud-blob/s3"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
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

		err := RegisterCache(ctx, scheme, NewBlobCache)

		if err != nil {
			return fmt.Errorf("Failed to register blob cache for '%s', %w", scheme, err)
		}

		register_map[scheme] = true
	}

	return nil
}

// NewBlobCache returns a BlobCache using the GoCloud package.
func NewBlobCache(ctx context.Context, uri string) (Cache, error) {

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

	r, err := bc.bucket.NewReader(ctx, uri, nil)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	return io.ReadAll(r)
}

// Set writes data to a BlobCache location.
func (bc *BlobCache) Set(uri string, body []byte) error {

	ctx := context.Background()

	var wr_opts *blob.WriterOptions

	// see notes above in NewBlobCacheFromURI

	logger := slog.Default()
	logger = logger.With("bucket uri", bc.bucket_uri)
	logger = logger.With("uri", uri)

	if strings.HasPrefix(bc.scheme, "s3") && bc.acl != "" {

		acl, err := aa_s3.StringACLToObjectCannedACL(bc.acl)

		if err != nil {
			logger.Error("Failed to derive ACL object", "error", err)
			return fmt.Errorf("Failed to derive ACL object, %w", err)
		}

		before := func(asFunc func(interface{}) bool) error {

			req := &aws_s3.PutObjectInput{}
			ok := asFunc(&req)

			if !ok {
				logger.Error("Invalid S3 type (asFunc)")
				return errors.New("invalid s3 type")
			}

			logger.Debug("Set ACL", "acl", acl)
			req.ACL = acl

			return nil
		}

		wr_opts = &blob.WriterOptions{
			BeforeWrite: before,
		}
	}

	wr, err := bc.bucket.NewWriter(ctx, uri, wr_opts)

	if err != nil {
		logger.Error("Failed to create new blob writer", "error", err)
		return err
	}

	_, err = wr.Write(body)

	if err != nil {
		logger.Error("Failed to write blob", "error", err)
		wr.Close()
		return err
	}

	err = wr.Close()

	if err != nil {
		logger.Error("Failed to close blob", "error", err)
		return err
	}

	logger.Debug("Successfully wrote blob")
	return nil
}

// Unset deletes data from a BlobCache location.
func (bc *BlobCache) Unset(uri string) error {
	ctx := context.Background()
	return bc.bucket.Delete(ctx, uri)
}

func (bc *BlobCache) Close() error {
	return bc.bucket.Close()
}
