# go-aws-auth

Go package providing methods and tools for determining or assigning AWS credentials.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-aws-auth.svg)](https://pkg.go.dev/github.com/aaronland/go-aws-auth)

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/aws-sts-session cmd/aws-sts-session/main.go
go build -mod vendor -ldflags="-s -w" -o bin/aws-mfa-session cmd/aws-mfa-session/main.go
go build -mod vendor -ldflags="-s -w" -o bin/aws-get-credentials cmd/aws-get-credentials/main.go
go build -mod vendor -ldflags="-s -w" -o bin/aws-cognito-credentials cmd/aws-cognito-credentials/main.go
go build -mod vendor -ldflags="-s -w" -o bin/aws-set-env cmd/aws-set-env/main.go
go build -mod vendor -ldflags="-s -w" -o bin/aws-sign-request cmd/aws-sign-request/main.go
go build -mod vendor -ldflags="-s -w" -o bin/aws-credentials-json-to-ini cmd/aws-credentials-json-to-ini/main.go
go build -mod vendor -ldflags="-s -w" -o bin/aws-imds-credentials cmd/aws-imds-credentials/main.go
```

## aws-cognito-credentials

`aws-cognito-credentials` generates temporary STS credentials for a given user in a Cognito identity pool.

```
$> ./bin/aws-cognito-credentials -h
Usage of ./bin/aws-cognito-credentials:
  -aws-config-uri string
    	A valid github.com/aaronland/go-aws-auth.Config URI.
  -duration int
    	The duration, in seconds, of the role session. Can not be less than 900. (default 900)
  -identity-pool-id string
    	A valid AWS Cognito Identity Pool ID.
  -login value
    	One or more key=value strings mapping to AWS Cognito authentication providers.
  -role-arn string
    	A valid AWS IAM role ARN to assign to STS credentials.
  -role-session-name string
    	An identifier for the assumed role session.
  -session-policy value
    	Zero or more IAM ARNs to use as session policies to supplement the default role ARN.	
```

For example:

```
$> go bin/aws-cognito-credentials \
	-aws-config-uri 'aws://us-east-1?credentials=session' \
	-identity-pool-id us-east-1:{GUID} \
	-login org.sfomuseum=bob
	-role-session-name bob -role-arn 'arn:aws:iam::{ACCOUNT_ID}:role/{ROLE}' \
	
| jq
	
{
  "AccessKeyId": "...",
  "Expiration": "...",
  "SecretAccessKey": "...",
  "SessionToken": "..."
}
```

### aws-credentials-json-to-ini

`aws-credentials-json-to-ini` reads JSON-encoded AWS credentials information and generates an AWS ini-style configuration file with those data.

```
$> ./bin/aws-credentials-json-to-ini -h
Usage of ./bin/aws-credentials-json-to-ini:
  -ini string
    	Path to the ini-style file where AWS credentials should be written. If "-" then data will be written to STDOUT.
  -json string
    	Path to the JSON file containing AWS credentials. If "-" then data will be read from STDIN.
  -name string
    	The name of the ini section where AWS credentials should be written. (default "default")
  -region string
    	The AWS region for the AWS credentials. (default "us-east-1")
```

For example:

```
$> go bin/aws-cognito-credentials \
	-aws-config-uri 'aws://us-east-1?credentials=session' \
	-identity-pool-id us-east-1:{GUID} \
	-login org.sfomuseum=bob
	-role-session-name bob -role-arn 'arn:aws:iam::{ACCOUNT_ID}:role/{ROLE}' \

| ./bin/aws-credentials-json-to-ini -json - -ini -

[default]
region = us-east-1
aws_access_key_id = ...
aws_secret_access_key = ...
aws_session_token = ...
```

### aws-get-credentials

`aws-get-credentials` is a command line tool to emit one or more keys from a given profile in an AWS .credentials file.

```
$> ./bin/aws-get-credentials -h
Usage of ./bin/aws-get-credentials:
  -profile string
    	A valid AWS credentials profile (default "default")
```

### aws-imds-credentials

`aws-imds-credentials` returns the current `aws.Credentials` derived from the EC2 IMDS API. For example:

```
$> ./bin/aws-imds-credentials | jq
{
  "AccessKeyID": "...",
  "SecretAccessKey": "...",
  "SessionToken": "...",
  "Source": "EC2RoleProvider",
  "CanExpire": true,
  "Expires": "2024-03-28T19:44:42.59621653Z"
}
```

### aws-mfa-session

`aws-mfa-session` is a command line to create session-based authentication keys and secrets for a given profile and multi-factor authentication (MFA) token and then writing that key and secret back to a "credentials" file in a specific profile section.

```
$> ./bin/aws-mfa-session -h
Usage of ./bin/aws-mfa-session:
  -duration string
    	A valid ISO8601 duration string indicating how long the session should last (months are currently not supported) (default "PT1H")
  -profile string
    	A valid AWS credentials profile (default "default")
  -session-profile string
    	The name of the AWS credentials profile to update with session credentials (default "session")
```

For example:

```
$> ./bin/aws-mfa-session -profile {PROFILE} -duration PT8H
Enter your MFA token code: 123456
2018/07/26 09:47:09 Updated session credentials for 'session' profile, expires Jul 26 17:47:09 (2018-07-27 00:51:52 +0000 UTC)
```

### aws-set-env

`aws-set-env` is a command line tool to assign required AWS authentication environment variables for a given profile in a AWS .credentials file.

```
$> ./bin/aws-set-env -h
Usage of ./bin/aws-set-env:
  -profile string
    	A valid AWS credentials profile (default "default")
  -session-token
    	Require AWS_SESSION_TOKEN environment variable (default true)
```

### aws-sign-request

`aws-sign-request` signs a HTTP request with an AWS "v4" signature, optionally executing the request and emitting the output to STDOUT or writing the request itself to STDOUT.

```
$> ./bin/aws-sign-request -h
Usage of ./bin/aws-sign-request:
  -api-signing-name string
    	The name the API uses to identify the service the request is scoped to.
  -api-signing-region string
    	If empty then the value of the region associated with the AWS config/credentials will be used.
  -credentials-uri string
    	A valid aaronland/go-aws-auth config URI.
  -debug
    	Enable verbose debug logging to STDOUT.	
  -do
    	If true then execute the signed request and output the response to STDOUT.
  -header value
    	Zero or more HTTP headers to assign to the request in the form of key=value.
  -method string
    	A valid HTTP method. (default "GET")
  -uri string
    	The URI you are trying to sign.
```

For example, to call a Lambda Function URL:

```
$> bin/aws-sign-request \
	-credentials-uri 'aws://{REGION}?credentials=iam:' \
	-api-signing-name 'lambda' \
	-uri https://{GIBBERISH}.lambda-url.{REGION}.on.aws/api/point-in-polygon \
	-method POST \
	-do \
	'{"latitude": 25.0, "longitude": -45.6 }' \
	
	| jq

{
  "places": [
    {
      "wof:id": "404528709",
      "wof:parent_id": "-1",
      "wof:name": "North Atlantic Ocean",
      "wof:country": "",
      "wof:placetype": "ocean",
      "mz:latitude": 0,
      "mz:longitude": 0,
      "mz:min_latitude": 24.965357,
      "mz:min_longitude": 0,
      "mz:max_latitude": -45.616087,
      "mz:max_longitude": -45.570425,
      "mz:is_current": 1,
      "mz:is_deprecated": -1,
      "mz:is_ceased": -1,
      "mz:is_superseded": 0,
      "mz:is_superseding": 0,
      "edtf:inception": "",
      "edtf:cessation": "",
      "wof:supersedes": [],
      "wof:superseded_by": [],
      "wof:belongsto": [],
      "wof:path": "404/528/709/404528709.geojson",
      "wof:repo": "whosonfirst-data-admin-xy",
      "wof:lastmodified": 1690923898
    }
  ]
}
```

### aws-sts-session

Generate STS credentials for a given profile and MFA token and then write those credentials back to an AWS "credentials" file in a specific profile section.

```
$> ./bin/aws-sts-session -h
Generate STS credentials for a given profile and MFA token and then write those credentials back to an AWS "credentials" file in a specific profile section.
Usage:
	 ./bin/aws-sts-session [options]
Valid options are:
  -config-uri string
    	A valid aaronland/gp-aws-auth.Config URI.
  -mfa
    	Require a valid MFA token code when assuming role. (default true)
  -mfa-serial-number string
    	The unique identifier of the MFA device being used for authentication.
  -mfa-token string
    	A valid MFA token string. If empty then data will be read from a command line prompt.
  -role-arn string
    	The AWS role ARN URI of the role you want to assume.
  -role-duration int
    	The duration, in seconds, of the role session. (default 3600)
  -role-session string
    	A unique name to identify the session.
  -session-profile string
    	The name of the AWS credentials profile to associate the temporary credentials with.
```

For example:

```
$> ./bin/aws-sts-session -config-uri 'aws://?region={REGION}&credentials={CREDENTIALS}' \
	-role-arn 'arn:aws:iam::{AWS_ACCOUNT}:role/{IAM_ROLE}' \
	-role-session debug \
	-mfa-serial-number arn:aws:iam::{AWS_ACCOUNT}:mfa/{MFA_LABEL} \
	-mfa-token {TOKEN} \
	-session-profile test

2024/11/08 08:23:25 Assumed role "arn:aws:sts::{AWS_ACCOUNT}:assumed-role/{IAM_ROLE}/debug", expires 2024-11-08 17:23:25 +0000 UTC
```

Note that this assumes a role with a "trust policy" equivalent to this:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Statement1",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::{AWS_ACCOUNT}:user/{IAM_USER}"
            },
            "Action": "sts:AssumeRole",
            "Condition": {
                "Bool": {
                    "aws:MultiFactorAuthPresent": true
                }
            }
        }
    ]
}
```

## Credentials

Credentials for URIs are defined as string labels. They are:

| Label | Description |
| --- | --- |
| `anon:` | Empty or anonymous credentials. |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `sts:{ARN}` | Assume the role defined by `{ARN}` using STS credentials. |
| `{AWS_PROFILE_NAME}` | This this profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | This this profile from a user-defined AWS credentials location. |

For example:

```
aws:///us-east-1?credentials=iam:
```

## See also:

* https://github.com/aws/aws-sdk-go-v2/