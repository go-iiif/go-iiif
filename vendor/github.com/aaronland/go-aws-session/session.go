package session

import (
	"github.com/aaronland/go-string/dsn"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
)

func NewSessionWithDSN(dsn_str string) (*aws_session.Session, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(dsn_str, "credentials", "region")

	if err != nil {
		return nil, err
	}

	return NewSessionWithCredentials(dsn_map["credentials"], dsn_map["region"])
}

func NewSessionWithCredentials(str_creds string, region string) (*aws_session.Session, error) {

	cfg, err := NewConfigWithCredentialsAndRegion(str_creds, region)

	if err != nil {
		return nil, err
	}

	sess := aws_session.New(cfg)

	_, err = sess.Config.Credentials.Get()

	if err != nil {
		return nil, err
	}

	return sess, nil
}
