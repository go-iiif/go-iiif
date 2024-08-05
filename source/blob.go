package source

import (
	"context"
	_ "log"

	"github.com/aaronland/gocloud-blob/bucket"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	"gocloud.dev/blob"
)

type BlobSource struct {
	Source
	bucket *blob.Bucket
}

func NewBlobSourceURIFromConfig(config *iiifconfig.Config) (string, error) {

	uri := cfg.Source.URI

	if uri == "" {
		cfg := config.Images
		uri = cfg.Source.Path
	}

	return uri, nil
}

func NewBlobSource(config *iiifconfig.Config) (Source, error) {

	uri, err := NewBlobSourceURIFromConfig(config)

	if err != nil {
		return nil, err
	}

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
	return bs.bucket.ReadAll(ctx, uri)
}
