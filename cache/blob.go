package cache

import (
	"context"
	"errors"
	"github.com/aaronland/gocloud-blob-bucket"
	"github.com/aws/aws-sdk-go/aws"	
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	"gocloud.dev/blob"
	"io/ioutil"
	"log"
	"net/url"
)

type BlobCache struct {
	Cache
	bucket_uri string
	bucket     *blob.Bucket
	scheme     string
	acl        string
}

func NewBlobCache(config iiifconfig.CacheConfig) (Cache, error) {

	uri := config.Path

	return NewBlobCacheFromURI(uri)
}

func NewBlobCacheFromURI(uri string) (Cache, error) {

	ctx := context.Background()

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(uri)
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

func (bc *BlobCache) Exists(uri string) bool {

	ctx := context.Background()

	exists, err := bc.bucket.Exists(ctx, uri)

	if err != nil {
		return false
	}

	return exists
}

func (bc *BlobCache) Get(uri string) ([]byte, error) {

	ctx := context.Background()

	fh, err := bc.bucket.NewReader(ctx, uri, nil)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return ioutil.ReadAll(fh)
}

// something something something permissions and ACLs
// https://github.com/google/go-cloud/issues/1108
// https://godoc.org/gocloud.dev/blob

func (bc *BlobCache) Set(uri string, body []byte) error {

	ctx := context.Background()

	bc.acl = "public-read"

	log.Println("SET", uri, bc.scheme, bc.acl)

	var wr_opts *blob.WriterOptions

	if bc.scheme == "s3" && bc.acl != "" {

		before := func(asFunc func(interface{}) bool) error {

			req := &s3manager.UploadInput{}
			ok := asFunc(&req)

			if !ok {
				return errors.New("invalid s3 type")
			}

			req.ACL = aws.String(bc.acl)
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

func (bc *BlobCache) Unset(uri string) error {
	ctx := context.Background()
	return bc.bucket.Delete(ctx, uri)
}
