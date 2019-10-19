package bucket

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronland/go-aws-session"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	"gocloud.dev/blob/s3blob"
	_ "log"
	"net/url"
	"strings"
)

func OpenBucket(ctx context.Context, uri string) (*blob.Bucket, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	query := u.Query()

	var bckt *blob.Bucket

	switch u.Scheme {
	case "s3":

		dsn := make([]string, 0)

		for k, v := range query {

			if len(v) != 1 {
				return nil, errors.New("Invalid DSN value")
			}

			dsn = append(dsn, fmt.Sprintf("%s=%s", k, v[0]))
		}

		str_dsn := strings.Join(dsn, " ")

		sess, err := session.NewSessionWithDSN(str_dsn)

		if err != nil {
			return nil, err
		}

		b, err := s3blob.OpenBucket(ctx, sess, u.Host, nil)

		if err != nil {
			return nil, err
		}

		bckt = b

	default:

		b, err := blob.OpenBucket(ctx, uri)

		if err != nil {
			return nil, err
		}

		bckt = b
	}

	// https://gocloud.dev/howto/blob/#prefix

	prefix := query.Get("prefix")

	if prefix != "" {

		if !strings.HasSuffix(prefix, "/") {
			prefix = fmt.Sprintf("%s/", prefix)
		}

		bckt = blob.PrefixedBucket(bckt, prefix)
	}

	return bckt, nil
}
