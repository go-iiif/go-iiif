package s3blob

import (
	"context"
	"fmt"
	_ "log/slog"
	"net/url"
	"sync"

	"github.com/aaronland/go-aws-auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gocloud.dev/blob"
	gc_s3blob "gocloud.dev/blob/s3blob"
)

const Scheme = "s3blob"

func init() {
	blob.DefaultURLMux().RegisterBucket(Scheme, new(lazySessionOpener))
}

type URLOpener struct {
	config *aws.Config
	bucket string
	prefix string
}

type lazySessionOpener struct {
	init   sync.Once
	opener *URLOpener
	err    error
}

func (o *lazySessionOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {

	o.init.Do(func() {

		bucket := u.Host

		q := u.Query()
		prefix := q.Get("prefix")

		auth_q := url.Values{}
		auth_q.Set("region", q.Get("region"))
		auth_q.Set("credentials", q.Get("credentials"))

		auth_uri := url.URL{}
		auth_uri.Scheme = "aws"
		auth_uri.RawQuery = auth_q.Encode()

		cfg, err := auth.NewConfig(ctx, auth_uri.String())

		if err != nil {
			o.err = err
			return
		}

		o.opener = &URLOpener{
			config: &cfg,
			bucket: bucket,
			prefix: prefix,
		}
	})

	if o.err != nil {
		return nil, fmt.Errorf("open bucket %v: %v", u, o.err)
	}

	return o.opener.OpenBucketURL(ctx, u)
}

func (o *URLOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {

	s3_client := s3.NewFromConfig(*o.config)
	s3_bucket := o.bucket
	s3_prefix := o.prefix

	b, err := gc_s3blob.OpenBucketV2(ctx, s3_client, s3_bucket, nil)

	if err != nil {
		return nil, err
	}

	if s3_prefix != "" {
		b = blob.PrefixedBucket(b, s3_prefix)
	}

	return b, nil
}
