package source

import (
	"context"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	"gocloud.dev/blob"
	_ "log"
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

	b, err := blob.OpenBucket(ctx, uri)

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
	return bs.bucket.ReadAll(ctx, uri)
}
