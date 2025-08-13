package s3

import (
	"context"
	"fmt"

	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"gocloud.dev/blob"
)

// NewWriterWithACL returns a new `blob.Writer` instance that has been configured with the relevant
// `blob.WriterOptions` to ensure that files written to S3 will be done using AWS ACL permissions
// defined in 'acl'.
func NewWriterWithACL(ctx context.Context, bucket *blob.Bucket, path string, str_acl string) (*blob.Writer, error) {

	acl, err := StringACLToObjectCannedACL(str_acl)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive canned ACL from string, %w", err)
	}

	before := func(asFunc func(interface{}) bool) error {

		req := &aws_s3.PutObjectInput{}
		ok := asFunc(&req)

		if ok {
			req.ACL = acl
		}

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
