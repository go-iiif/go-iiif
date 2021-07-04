# gocloud-blob-s3

This is thing wrapper around the default `go-cloud` S3 blob opener to check for a `credentials` parameter (in blob URIs) and use it to assign AWS S3 session credentials.

## Example

```			
import (
	"context"
	"gocloud.dev/blob"
	_ "github.com/aaronland/gocloud-blob-s3"
)

func main() {

	ctx := context.Background()	
	bucket, _ := blob.OpenBucket(ctx, "s3blob://BUCKET?region=REGION&credentials=CREDENTIALS")

	// do stuff with bucket here
}
```

_Note the use of the `s3blob://` scheme which is different than the default `s3://` scheme._

## Credentials

Credentials for AWS sessions are defined as string labels. They are:

| Label | Description |
| --- | --- |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `{AWS_PROFILE_NAME}` | This this profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | This this profile from a user-defined AWS credentials location. |

## See also

* https://github.com/aaronland/go-aws-session
* https://gocloud.dev/howto/blob