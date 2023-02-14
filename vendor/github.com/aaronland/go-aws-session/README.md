# go-aws-session

This package is a thing wrapper around the [AWS Go SDK](https://docs.aws.amazon.com/sdk-for-go) to allow for creating sessions using DSN strings.

## Example

```

import (
	"github.com/aaronland/go-aws-session"
)

func main() {
	str_dsn := "region=us-east-1 credentials=env:"
	sess, err := session.NewSessionWithDSN(str_dsn)

	// do something with sess or err here
}

```

## DSN strings

The following properties are required in DSN strings:

### Credentials

Credentials for AWS sessions are defined as string labels. They are:

| Label | Description |
| --- | --- |
| `anon:` | Empty or anonymous credentials. |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `{AWS_PROFILE_NAME}` | This this profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | This this profile from a user-defined AWS credentials location. |

For example:

```
s3:///bucket-name?region=us-east-1&credentials=iam:
```

### Region

Any valid AWS region.

## See also

* https://aws.amazon.com/blogs/security/a-new-and-standardized-way-to-manage-credentials-in-the-aws-sdks/
* https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html
* https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html
* https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
* https://github.com/google/go-cloud/blob/master/aws/aws.go
