package s3blob

import (
	"context"
)

// SetACLWriterOptionsWithContext return a new context.Context instance with a gocloud.dev/blob.WriterOptions
// instance used to assign 'acl' permissions for all S3 blob writes. The WriterOptions instance is assigned
// to the new context with key 'key' and is assumed to be retrieved later by code using blob.NewWriter instances.
// This method is DEPRECATED. Please use SetWriterOptionsWithContext() instead.
func SetACLWriterOptionsWithContext(ctx context.Context, key interface{}, acl string) context.Context {
	ctx, _ = SetWriterOptionsWithContext(ctx, key, "ACL", acl)
	return ctx
}
