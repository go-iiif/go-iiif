package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os/user"
	"path/filepath"
	_ "sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// null_cfg is a placeholder to return in error contexts.
var null_cfg aws.Config

const AnonymousCredentialsString string = "anon:"

const EnvironmentCredentialsString string = "env:"

const IAMCredentialsString string = "iam:"

const ProfileCredentialsTemplate string = "{profile}"

const ProfileFromPathCredentialsTemplate string = "{path}:{profile}"

const StaticCredentialsTemplate string = "static:{id}:{key}:{secret}"

const StaticCredentialsPrefix string = "static:"

const STSCredentialsPrefix string = "sts:"

// ValidCredentials() returns the list of valid credential strings patterns.
func ValidCredentials() []string {

	valid := []string{
		AnonymousCredentialsString,
		EnvironmentCredentialsString,
		IAMCredentialsString,
		STSCredentialsPrefix,
		ProfileCredentialsTemplate,
		ProfileFromPathCredentialsTemplate,
		StaticCredentialsTemplate,
	}

	return valid
}

func ValidCredentialsString() string {

	valid := ValidCredentials()

	return fmt.Sprintf("Valid credential flags are: %s", strings.Join(valid, ", "))
}

// CredentialsStringPatterns() returns the list of valid credential strings patterns. DEPRECATED. Please use `ValidCredentials` instead.
func CredentialsStringPatterns() []string {
	slog.Warn("CredentialsStringPatterns() is deprecated. Please call ValidCredentials() instead.")
	return ValidCredentials()
}

// NewConfig() returns a new `aws.Config` derived from 'uri' which is expected to be configured
// in the form of:
//
//	aws://{AWS_REGION}?credentials={CREDENTIALS_STRING}
//
// Where {AWS_REGION} is a valid AWS region name and {CREDENTIALS_STRING} is a string in the form of:
//
//	`anon:` Use anonymous credentials
//	`env:` Use credentials derived from "AWS_" environment variables
//	`iam:` Use IAM credentials
//	`{PROFILE}` Use a specific profile defined by {PROFILE} from the default credentials file
//	`{PATH}:{PROFILE}` Use a specific profile definied by {PROFILE} from the credentials file defined by {PATH}
//	`static:{KEY}:{SECRET}:{TOKEN}` Read credentials as positional elements in a string
//	`` If credentials are passed as an empty string then use default credentials strategy defined by `aws-sdk-go-v2`
func NewConfig(ctx context.Context, uri string) (aws.Config, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return null_cfg, fmt.Errorf("Failed to parse URI, %w", err)
	}

	region := u.Host

	q := u.Query()

	creds := q.Get("credentials")

	if region == "" {
		region = q.Get("region")
	}

	cfg, err := NewConfigWithCredentialsString(ctx, creds)

	if err != nil {
		return null_cfg, fmt.Errorf("Failed to derive config from credentials string, %w", err)
	}

	cfg.Region = region
	return cfg, nil
}

// NewConfigWithCredentialsString() returns a new `aws.Config` derived from 'str_creds' which is expected to be passed
// in as one of the following:
//
//	`anon:` Use anonymous credentials
//	`env:` Use credentials derived from "AWS_" environment variables
//	`iam:` Use IAM credentials
//	`{PROFILE}` Use a specific profile defined by {PROFILE} from the default credentials file
//	`{PATH}:{PROFILE}` Use a specific profile definied by {PROFILE} from the credentials file defined by {PATH}
//	`static:{KEY}:{SECRET}:{TOKEN}` Read credentials as positional elements in a string
//	`` If credentials are passed as an empty string then use default credentials strategy defined by `aws-sdk-go-v2`
func NewConfigWithCredentialsString(ctx context.Context, str_creds string) (aws.Config, error) {

	if strings.HasPrefix(str_creds, "anon:") {

		provider := aws.AnonymousCredentials{}

		return config.LoadDefaultConfig(ctx,
			config.WithCredentialsProvider(provider),
		)

	} else if strings.HasPrefix(str_creds, "static:") {

		details := strings.Split(str_creds, ":")

		if len(details) != 4 {
			return null_cfg, fmt.Errorf("Expected (4) components for 'static:' credentials URI but got %d", len(details))
		}

		key := details[1]
		secret := details[2]
		token := details[3]

		provider := credentials.NewStaticCredentialsProvider(key, secret, token)

		return config.LoadDefaultConfig(ctx,
			config.WithCredentialsProvider(provider),
		)

	} else if str_creds == "iam:" || str_creds == "env:" {

		return config.LoadDefaultConfig(ctx)

	} else if str_creds != "" {

		details := strings.Split(str_creds, ":")

		var creds_file string
		var profile string

		if len(details) == 1 {

			whoami, err := user.Current()

			if err != nil {
				return null_cfg, fmt.Errorf("Failed to determine current user, %w", err)
			}

			dotaws := filepath.Join(whoami.HomeDir, ".aws")
			creds_file = filepath.Join(dotaws, "credentials")

			profile = details[0]

		} else {

			path, err := filepath.Abs(details[0])

			if err != nil {
				return null_cfg, fmt.Errorf("Failed to derive absolute path for %s, %w", details[0], err)
			}

			creds_file = path
			profile = details[1]
		}

		return config.LoadDefaultConfig(ctx,
			config.WithSharedCredentialsFiles([]string{creds_file}),
			config.WithSharedConfigProfile(profile),
		)

	} else {

		return null_cfg, fmt.Errorf("Invalid or unsupported credentials string")
	}

}
