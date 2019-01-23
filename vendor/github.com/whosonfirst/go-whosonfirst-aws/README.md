# go-whosonfirst-aws

There are many AWS wrappers. This one is ours.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

This works. Until it doesn't. It has not been properly documented yet.

## DSN strings

```
bucket=BUCKET region={REGION} prefix={PREFIX} credentials={CREDENTIALS}
```

Valid credentials strings are:

* `env:`

* `iam:`

* `{PATH}:{PROFILE}`

## See also

* https://docs.aws.amazon.com/sdk-for-go/

