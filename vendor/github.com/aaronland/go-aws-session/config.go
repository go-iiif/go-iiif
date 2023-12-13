package session

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewConfigWithCredentialsAndRegion(str_creds string, region string) (*aws.Config, error) {

	cfg, err := NewConfigWithCredentials(str_creds)

	if err != nil {
		return nil, fmt.Errorf("Failed to create config from credentials, %w", err)
	}

	cfg.WithRegion(region)
	return cfg, nil
}

func NewConfigWithCredentials(str_creds string) (*aws.Config, error) {

	cfg := aws.NewConfig()

	if strings.HasPrefix(str_creds, AnonymousCredentialsString) {

		creds := credentials.AnonymousCredentials
		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(str_creds, EnvironmentCredentialsString) {

		creds := credentials.NewEnvCredentials()
		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(str_creds, STSCredentialsPrefix) {

		// https://github.com/aws/aws-sdk-go/issues/801
		// https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/stscreds/

		sess, err := session.NewSession()

		if err != nil {
			return nil, fmt.Errorf("Failed to create new session for %s provider, %w", STSCredentialsPrefix, err)
		}

		arn := strings.Replace(str_creds, STSCredentialsPrefix, "", 1)

		session_name := filepath.Base(arn)

		creds := stscreds.NewCredentials(sess, arn, func(provider *stscreds.AssumeRoleProvider) {
			provider.RoleARN = arn
			provider.RoleSessionName = session_name
		})

		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(str_creds, IAMCredentialsString) {

		// Do nothing...

	} else if strings.HasPrefix(str_creds, StaticCredentialsPrefix) {

		details := strings.Split(str_creds, ":")

		if len(details) != 4 {
			return nil, fmt.Errorf("Expected (4) components for 'static:' credentials URI but got %d", len(details))
		}

		id := details[1]
		key := details[2]
		secret := details[3]

		creds := credentials.NewStaticCredentials(id, key, secret)
		cfg.WithCredentials(creds)

	} else if str_creds != "" {

		details := strings.Split(str_creds, ":")

		var creds_file string
		var profile string

		if len(details) == 1 {

			whoami, err := user.Current()

			if err != nil {
				return nil, fmt.Errorf("Failed to derive current user, %w", err)
			}

			dotaws := filepath.Join(whoami.HomeDir, ".aws")
			creds_file = filepath.Join(dotaws, "credentials")

			profile = details[0]

		} else {

			path, err := filepath.Abs(details[0])

			if err != nil {
				return nil, fmt.Errorf("Failed to derive absolute path for '%s', %w", details[0], err)
			}

			creds_file = path
			profile = details[1]
		}

		creds := credentials.NewSharedCredentials(creds_file, profile)
		cfg.WithCredentials(creds)

	} else {

		// for backwards compatibility as of 05a6042dc5956c13513bdc5ab4969877013f795c
		// (20161203/thisisaaronland)

		creds := credentials.NewEnvCredentials()
		cfg.WithCredentials(creds)
	}

	return cfg, nil
}
