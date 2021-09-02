package s3blob

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gocloud.dev/blob"
)

type AsFunc func(func(interface{}) bool) error

// SetACLWriterOptionsWithContext return a new context.Context instance with a gocloud.dev/blob.WriterOptions
// instance whose properties are assigned according to 'opt_key' and 'opt_value', where the latter is a valid
// blob.WriterOptions property and the former is a valid type for that property. The one special-case is the
// "ACL" opt_key which which takes a opt_value of type type string and will assign a suitable GoCloud BeforeWrite
// option for assigning AWS S3 ACL permissions.
//
// The WriterOptions instance is assigned to the new context with key 'ctx_key' and is assumed to be retrieved later
// by code using blob.NewWriter instances.
func SetWriterOptionsWithContext(ctx context.Context, ctx_key interface{}, opt_key string, opt_value interface{}) (context.Context, error) {

	var wr_opts *blob.WriterOptions

	v := ctx.Value(ctx_key)

	if v == nil {
		wr_opts = &blob.WriterOptions{}
	} else {
		wr_opts = v.(*blob.WriterOptions)
	}

	switch opt_key {
	case "BufferSize":

		switch opt_value.(type) {
		case int:
			// pass
		default:
			return nil, fmt.Errorf("Invalid type for '%s' value", opt_key)
		}

		wr_opts.BufferSize = opt_value.(int)

	case "CacheControl", "ContentDisposition", "ContentEncoding", "ContentType", "ContentLanguage":

		switch opt_value.(type) {
		case string:
			// pass
		default:
			return nil, fmt.Errorf("Invalid type for '%s' value", opt_key)
		}

		switch opt_key {
		case "CacheControl":
			wr_opts.CacheControl = opt_value.(string)
		case "ContentDisposition":
			wr_opts.ContentDisposition = opt_value.(string)
		case "ContentEncoding":
			wr_opts.ContentEncoding = opt_value.(string)
		case "ContentType":
			wr_opts.ContentType = opt_value.(string)
		case "ContentLanguage":
			wr_opts.ContentLanguage = opt_value.(string)
		}

	case "ContentMD5":

		switch opt_value.(type) {
		case string:
			// pass
		default:
			return nil, fmt.Errorf("Invalid type for '%s' value", opt_key)
		}

		wr_opts.ContentMD5 = opt_value.([]byte)

	case "Metadata":

		switch opt_value.(type) {
		case map[string]string:
			// pass
		default:
			return nil, fmt.Errorf("Invalid type for '%s' value", opt_key)
		}

		wr_opts.Metadata = opt_value.(map[string]string)

	case "BeforeWrite":

		switch opt_value.(type) {
		case AsFunc:
			// pass
		default:
			return nil, fmt.Errorf("Invalid type for '%s' value", opt_key)
		}

		wr_opts.BeforeWrite = opt_value.(AsFunc)

	case "ACL":

		switch opt_value.(type) {
		case string:
			// pass
		default:
			return nil, fmt.Errorf("Invalid type for '%s' value", opt_key)
		}

		acl := opt_value.(string)

		before := func(asFunc func(interface{}) bool) error {

			req := &s3manager.UploadInput{}
			ok := asFunc(&req)

			if !ok {
				return fmt.Errorf("invalid s3 type")
			}

			req.ACL = aws.String(acl)
			return nil
		}

		wr_opts.BeforeWrite = before

	default:
		return nil, fmt.Errorf("Invalid or unsupported key '%s'", opt_key)
	}

	ctx = context.WithValue(ctx, ctx_key, wr_opts)
	return ctx, nil
}

// SetWriterOptionsWithContextAndMap is a convenience method for invoking SetWriterOptionsWithContext
// multiple times.
func SetWriterOptionsWithContextAndMap(ctx context.Context, ctx_key interface{}, opts map[string]interface{}) (context.Context, error) {

	var err error

	for k, v := range opts {

		ctx, err = SetWriterOptionsWithContext(ctx, ctx_key, k, v)

		if err != nil {
			return nil, fmt.Errorf("Failed to set writer options for '%s', %v", k, v)
		}
	}

	return ctx, nil
}
