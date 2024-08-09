package source

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aaronland/gocloud-blob/bucket"
	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
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

		err := RegisterSource(ctx, scheme, NewBlobSourceFromURI)

		if err != nil {
			return fmt.Errorf("Failed to register blob source for '%s', %w", scheme, err)
		}

		register_map[scheme] = true
	}

	return nil
}

func NewBlobSourceURIFromConfig(cfg *iiifconfig.Config) (string, error) {

	uri := cfg.Images.Source.URI

	if uri == "" {
		uri = cfg.Images.Source.Path
	}

	return uri, nil
}

func NewBlobSource(cfg *iiifconfig.Config) (Source, error) {

	uri, err := NewBlobSourceURIFromConfig(cfg)

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
