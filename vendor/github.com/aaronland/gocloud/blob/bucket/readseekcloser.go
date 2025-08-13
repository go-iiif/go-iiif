package bucket

import (
	"context"
	"fmt"
	"io"

	"github.com/whosonfirst/go-ioutil"
	"gocloud.dev/blob"
)

func NewReadSeekCloser(ctx context.Context, b *blob.Bucket, key string, opts *blob.ReaderOptions) (io.ReadSeekCloser, error) {

	r, err := b.NewReader(ctx, key, opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new reader for %s, %w", key, err)
	}

	rsc, err := ioutil.NewReadSeekCloser(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new read seek closer for %s, %w", key, err)
	}

	return rsc, nil
}
