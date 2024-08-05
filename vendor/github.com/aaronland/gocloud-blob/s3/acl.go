package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// StringACLToObjectCannedACL resolves a subset of the string values for S3 ACLs (those specific to objects) to
// their corresponding `github.com/aws/aws-sdk-go-v2/service/s3/types.ObjectCannedACL` instance.
func StringACLToObjectCannedACL(str_acl string) (types.ObjectCannedACL, error) {

	switch str_acl {
	case "private":
		return types.ObjectCannedACLPrivate, nil
	case "public-read":
		return types.ObjectCannedACLPublicRead, nil
	case "public-read-write":
		return types.ObjectCannedACLPublicReadWrite, nil
	case "authenticated-read":
		return types.ObjectCannedACLAuthenticatedRead, nil
	default:
		return "", fmt.Errorf("Invalid or unsupported ACL")
	}

}
