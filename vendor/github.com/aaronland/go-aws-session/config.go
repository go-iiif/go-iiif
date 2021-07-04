package session

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os/user"
	"path/filepath"
	"strings"
)

func NewConfigWithCredentialsAndRegion(str_creds string, region string) (*aws.Config, error) {

	cfg, err := NewConfigWithCredentials(str_creds)

	if err != nil {
		return nil, err
	}

	cfg.WithRegion(region)
	return cfg, nil
}

func NewConfigWithCredentials(str_creds string) (*aws.Config, error) {

	cfg := aws.NewConfig()

	if strings.HasPrefix(str_creds, "anon:") {

		creds := credentials.AnonymousCredentials
		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(str_creds, "env:") {

		creds := credentials.NewEnvCredentials()
		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(str_creds, "iam:") {

		// assume an IAM role suffient for doing whatever

	} else if strings.HasPrefix(str_creds, "static:") {

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
				return nil, err
			}

			dotaws := filepath.Join(whoami.HomeDir, ".aws")
			creds_file = filepath.Join(dotaws, "credentials")

			profile = details[0]

		} else {

			path, err := filepath.Abs(details[0])

			if err != nil {
				return nil, err
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
