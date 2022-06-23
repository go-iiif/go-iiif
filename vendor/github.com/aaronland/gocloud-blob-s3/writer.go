package s3blob

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gocloud.dev/blob"
)

// NewWriterWithACL returns a new `blob.Writer` instance that has been configured with the relevant
// `blob.WriterOptions` to ensure that files written to S3 will be done using AWS ACL permissions
// defined in 'acl'.
func NewWriterWithACL(ctx context.Context, bucket *blob.Bucket, path string, acl string) (*blob.Writer, error) {

	before := func(asFunc func(interface{}) bool) error {

		req := &s3manager.UploadInput{}
		ok := asFunc(&req)

		if !ok {
			return fmt.Errorf("invalid S3 type")
		}

		req.ACL = aws.String(acl)
		return nil
	}

	wr_opts := &blob.WriterOptions{
		BeforeWrite: before,
	}

	wr, err := bucket.NewWriter(ctx, path, wr_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create writer for %s, %w", path, err)
	}

	return wr, nil
}
