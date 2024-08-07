package source

import (
	"context"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/memblob"
)

type MemorySource struct {
	Source
	uri    string
	bucket *blob.Bucket
}

func NewMemorySource(body []byte) (Source, error) {

	ctx := context.Background()

	b, err := blob.OpenBucket(ctx, "mem://")

	if err != nil {
		return nil, err
	}

	uri := "memory://"

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

func (bs *MemorySource) String() string {
	return bs.uri
}

func (bs *MemorySource) Read(uri string) ([]byte, error) {

	ctx := context.Background()
	return bs.bucket.ReadAll(ctx, bs.uri)
}
