package source

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aaronland/gocloud-blob/bucket"
	"gocloud.dev/blob"
)

type BlobSource struct {
	Source
	bucket     *blob.Bucket
	bucket_uri string
}

func init() {
	ctx := context.Background()
	err := RegisterBlobSourceSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterBloblSourceSchemes will ...
func RegisterBlobSourceSchemes(ctx context.Context) error {

	register_mu.Lock()
	defer register_mu.Unlock()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		slog.Debug("Register blob source scheme", "scheme", scheme)

		_, exists := register_map[scheme]

		if exists {
			continue
		}

		err := RegisterSource(ctx, scheme, NewBlobSource)

		if err != nil {
			return fmt.Errorf("Failed to register blob source for '%s', %w", scheme, err)
		}

		register_map[scheme] = true
	}

	return nil
}

func NewBlobSource(ctx context.Context, uri string) (Source, error) {

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	bs := &BlobSource{
		bucket:     b,
		bucket_uri: uri,
	}

	return bs, nil
}

func (bs *BlobSource) String() string {
	return bs.bucket_uri
}

func (bs *BlobSource) Read(uri string) ([]byte, error) {

	ctx := context.Background()
	return bs.bucket.ReadAll(ctx, uri)
}

func (bs *BlobSource) Close() error {
	return bs.bucket.Close()
}
