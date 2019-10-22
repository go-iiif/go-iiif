package source

import (
	"context"
	"github.com/aaronland/gocloud-blob-bucket"
	"gocloud.dev/blob"
)

type MemorySource struct {
	Source
	bucket *blob.Bucket
}

func NewMemorySource(body []byte) (Source, error) {

	ctx := context.Background()

	b, err := bucket.OpenBucket(ctx, "mem://")

	if err != nil {
		return nil, err
	}

	err = b.WriteAll(ctx, "URI", body, nil)

	if err != nil {
		return nil, err
	}

	bs := &MemorySource{
		bucket: b,
	}

	return bs, nil
}

func (bs *MemorySource) Read(uri string) ([]byte, error) {

	ctx := context.Background()
	return bs.bucket.ReadAll(ctx, "URI")
}
