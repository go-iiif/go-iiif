package source

import (
	"context"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	"github.com/aaronland/gocloud-blob-bucket"
	"gocloud.dev/blob"
	"io/ioutil"
)

type BlobSource struct {
	Source
	bucket *blob.Bucket
}

func NewBlobSource(config *iiifconfig.Config) (Source, error) {

	cfg := config.Images
	ctx := context.Background()
	
	b, err := bucket.OpenBucket(ctx, cfg.Source.Path)

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
