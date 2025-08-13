package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/go-ini/ini"
)

// type Credentials defines a struct for encapsulating information about an .ini style AWS credentials file
type Credentials struct {
	// Path is the path to the credentials file on disk
	Path string
	// ini is the `ini.File` instance derived from the credentials file
	ini *ini.File
}

// NewCredentials will iterate through the list of default shared credentials files locations and return a
// new `Credentials` instance for the first match that can be successfully opened and parsed as a .ini-style
// config file.
func NewCredentials() (*Credentials, error) {

	var ini_config *ini.File
	var ini_path string

	for _, path := range config.DefaultSharedCredentialsFiles {

		_, err := os.Stat(path)

		if err != nil {
			continue
		}

		i, err := ini.Load(path)

		if err != nil {
			return nil, fmt.Errorf("Failed to load config file at %s, %w", path, err)
		}

		ini_path = path
		ini_config = i
		break
	}

	if ini_config == nil {
		return nil, fmt.Errorf("Unable to load config file")
	}

	c := Credentials{
		Path: ini_path,
		ini:  ini_config,
	}

	return &c, nil
}

// AWSConfigWithProfile will return a new `aws.Config` instance for the profile matching 'profile'
// in 'c'. This method also assumes that AWS region information is stored in a "region" key inside
// a "default" section block of the underlying AWS credentials file.
func (c *Credentials) AWSConfigWithProfile(ctx context.Context, profile string) (aws.Config, error) {

	sect := c.ini.Section("default")
	region := sect.Key("region")

	return config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region.String()),
	)
}

// SetSessionCredentialsWithProfile will write AWS credentials information derived from 'creds' to a
// block labeled after 'profile' to the the underlying AWS credentials file.
func (c *Credentials) SetSessionCredentialsWithProfile(ctx context.Context, profile string, creds *types.Credentials) error {

	sect := c.ini.Section(profile)

	sect.Key("aws_access_key_id").SetValue(*creds.AccessKeyId)
	sect.Key("aws_secret_access_key").SetValue(*creds.SecretAccessKey)
	sect.Key("aws_session_token").SetValue(*creds.SessionToken)

	err := c.ini.SaveTo(c.Path)

	if err != nil {
		return fmt.Errorf("Failed to save date to credentials files, %w", err)
	}

	return nil
}
