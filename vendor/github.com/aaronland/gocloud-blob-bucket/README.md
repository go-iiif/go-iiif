# gocloud-blob-bucket

This package is a thin wrapper around the [Go Cloud](https://gocloud.dev/howto/blob/) package in order to be able to specify AWS S3 credentials using string values. Those values are:

* `env:` – read credentials from AWS defined environment variables.
* `iam:` – assume AWS IAM credentials).
* `{AWS_PROFILE_NAME}`.
* `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}`.

For example:

```
s3:///bucket-name?region=us-east-1&credentials=iam:
```

Additionally, if a Go Cloud URI contains a `prefix=` query parameter this package will automatically return a `blob.PrefixedBucket`.

## See also

* https://gocloud.dev/howto/blob/
* https://github.com/aaronland/go-aws-session