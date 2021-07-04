package s3blob

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"	
)

// SetACLWriterOptionsWithContext return a new context.Context instance with a gocloud.dev/blob.WriterOptions
// instance used to assign 'acl' permissions for all S3 blob writes. The WriterOptions instance is assigned
// to the new context with key 'key' and is assumed to be retrieved later by code using blob.NewWriter instances.
func SetACLWriterOptionsWithContext(ctx context.Context, key interface{}, acl string) context.Context {

	before := func(asFunc func(interface{}) bool) error {
		
		req := &s3manager.UploadInput{}
		ok := asFunc(&req)
		
		if !ok {
			return fmt.Errorf("invalid s3 type")
		}
		
		req.ACL = aws.String(acl)
		return nil
	}
	
	wr_opts := &blob.WriterOptions{
		BeforeWrite: before,
	}
	
	return context.WithValue(ctx, key, wr_opts)
}
