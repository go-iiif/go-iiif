package source

import (
	"context"
	"github.com/aaronland/gocloud-blob-bucket"
	"gocloud.dev/blob"
)

type MemorySource struct {
	Source
	uri    string
	bucket *blob.Bucket
}

func NewMemorySource(body []byte) (Source, error) {

	ctx := context.Background()

	b, err := bucket.OpenBucket(ctx, "mem://")

	if err != nil {
		return nil, err
	}

	uri := "mem"

	err = b.WriteAll(ctx, uri, body, nil)

	if err != nil {
		return nil, err
	}

	bs := &MemorySource{
		uri:    uri,
		bucket: b,
	}

	return bs, nil
}

func (bs *MemorySource) Read(uri string) ([]byte, error) {

	ctx := context.Background()
	return bs.bucket.ReadAll(ctx, bs.uri)
}
