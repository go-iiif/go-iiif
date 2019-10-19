package source

import (
	"context"
	"github.com/aaronland/gocloud-blob-bucket"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	"gocloud.dev/blob"
	"io/ioutil"
)

type BlobSource struct {
	Source
	bucket *blob.Bucket
}

func NewBlobSource(config *iiifconfig.Config) (Source, error) {

	cfg := config.Images
	uri := cfg.Source.Path

	return NewBlobSourceFromURI(uri)
}

func NewBlobSourceFromURI(uri string) (Source, error) {

	ctx := context.Background()

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	bs := &BlobSource{
		bucket: b,
	}

	return bs, nil
}

func (bs *BlobSource) Read(uri string) ([]byte, error) {

	ctx := context.Background()

	fh, err := bs.bucket.NewReader(ctx, uri, nil)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return ioutil.ReadAll(fh)
}
