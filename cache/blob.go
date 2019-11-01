package cache

import (
	"context"
	"github.com/aaronland/gocloud-blob-bucket"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	"gocloud.dev/blob"
	"io/ioutil"
	"log"
)

type BlobCache struct {
	Cache
	bucket_uri string
	bucket     *blob.Bucket
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

	bc := &BlobCache{
		bucket_uri: uri,
		bucket:     b,
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

func (bc *BlobCache) Set(uri string, body []byte) error {

	log.Println("BLOB", "SET", uri)
	
	ctx := context.Background()

	fh, err := bc.bucket.NewWriter(ctx, uri, nil)

	if err != nil {
		return err
	}

	_, err = fh.Write(body)

	if err != nil {
		fh.Close()
		return err
	}

	return fh.Close()
}

func (bc *BlobCache) Unset(uri string) error {
	ctx := context.Background()
	return bc.bucket.Delete(ctx, uri)
}
