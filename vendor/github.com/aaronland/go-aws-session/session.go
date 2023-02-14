package session

import (
	"fmt"
	"github.com/aaronland/go-string/dsn"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"net/url"
)

func NewSession(uri string) (*aws_session.Session, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	q := u.Query()

	creds := q.Get("credentials")
	region := q.Get("region")

	if creds == "" {
		return nil, fmt.Errorf("Missing ?credentials parameter")
	}

	if region == "" {
		return nil, fmt.Errorf("Missing ?region parameter")
	}

	return NewSessionWithCredentials(creds, region)
}

func NewSessionWithDSN(dsn_str string) (*aws_session.Session, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(dsn_str, "credentials", "region")

	if err != nil {
		return nil, fmt.Errorf("Failed to create DSN from string, %w", err)
	}

	return NewSessionWithCredentials(dsn_map["credentials"], dsn_map["region"])
}

func NewSessionWithCredentials(str_creds string, region string) (*aws_session.Session, error) {

	cfg, err := NewConfigWithCredentialsAndRegion(str_creds, region)

	if err != nil {
		return nil, fmt.Errorf("Failed to create config, %w", err)
	}

	sess := aws_session.New(cfg)

	if str_creds != AnonymousCredentialsString {

		_, err = sess.Config.Credentials.Get()

		if err != nil {
			return nil, fmt.Errorf("Failed to derive credentials from session, %w", err)
		}
	}

	return sess, nil
}
