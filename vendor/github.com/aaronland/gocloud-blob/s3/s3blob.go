package s3

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aaronland/go-aws-auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"gocloud.dev/blob"
	gc_s3blob "gocloud.dev/blob/s3blob"
)

const Scheme = "s3blob"

func init() {
	is_registered := false

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		if scheme == Scheme {
			is_registered = true
			break
		}
	}

	if !is_registered {
		blob.DefaultURLMux().RegisterBucket(Scheme, new(sessionOpener))
	}
}

type URLOpener struct {
	config *aws.Config
	bucket string
	prefix string
}

type sessionOpener struct {
	opener *URLOpener
	err    error
}

func (o *sessionOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {

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
		return nil, fmt.Errorf("Failed to create new AWS (auth) config, %w", err)
	}

	opener := &URLOpener{
		config: &cfg,
		bucket: bucket,
		prefix: prefix,
	}

	return opener.OpenBucketURL(ctx, u)
}

func (o *URLOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {

	s3_client := aws_s3.NewFromConfig(*o.config)
	s3_bucket := o.bucket
	s3_prefix := o.prefix

	b, err := gc_s3blob.OpenBucketV2(ctx, s3_client, s3_bucket, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to open bucket, %w", err)
	}

	if s3_prefix != "" {
		b = blob.PrefixedBucket(b, s3_prefix)
	}

	return b, nil
}
