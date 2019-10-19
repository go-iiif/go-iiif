package session

import (
	"github.com/aaronland/go-string/dsn"
	"os"
	"strings"
)

func ConfigureEnvWithDSN(aws_dsn string) error {

	aws_config, err := dsn.StringToDSNWithKeys(aws_dsn, "credentials", "region")

	if err != nil {
		return err
	}

	aws_region := aws_config["region"]

	err = assignEnv("AWS_DEFAULT_REGION", aws_region)

	if err != nil {
		return err
	}

	aws_creds := aws_config["credentials"]

	if strings.HasPrefix(aws_creds, "env:") {

		// assume environment variables already

	} else if strings.HasPrefix(aws_creds, "iam:") {

		// assume an IAM role suffient for doing whatever

	} else if aws_creds != "" {

		details := strings.SplitN(aws_creds, ":", 2)

		var profile string

		if len(details) == 1 {

			profile = details[0]

		} else {

			path := details[0]

			err := assignEnv("AWS_SHARED_CREDENTIALS_FILE", path)

			if err != nil {
				return err
			}

			profile = details[1]
		}

		err = assignEnv("AWS_DEFAULT_PROFILE", profile)

		if err != nil {
			return err
		}

	} else {

		// assume environment variables already
	}

	return nil
}

func assignEnv(aws_var string, value string) error {

	err := os.Setenv(aws_var, value)

	if err != nil {
		return &AssignEnvError{err, aws_var, value}
	}

	return nil
}
