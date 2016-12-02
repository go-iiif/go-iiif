package aws

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "log"
	"path/filepath"
)

type S3Connection struct {
	service *s3.S3
	bucket  string
	prefix  string
}

type S3Config struct {
	Bucket      string
	Prefix      string
	Region      string
	Credentials string	// see notes below
}

func NewS3Connection(s3cfg S3Config) (*S3Connection, error) {

	// https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

	cfg := aws.NewConfig()
	cfg.WithRegion(s3cfg.Region)

	if s3cfg.Credentials != "" {

		// to do: update this to allow a profile
		// to do: update this to allow other aws-sdk credential types
		// https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/
		// (20161202/thisisaaronland)
		
		creds := credentials.NewSharedCredentials(s3cfg.Credentials, "")
		cfg.WithCredentials(creds)
	}

	sess := session.New(cfg)

	if s3cfg.Credentials != "" {

		_, err := sess.Config.Credentials.Get()

		if err != nil {
			return nil, err
		}
	}

	service := s3.New(sess)

	c := S3Connection{
		service: service,
		bucket:  s3cfg.Bucket,
		prefix:  s3cfg.Prefix,
	}

	return &c, nil
}

func (conn *S3Connection) Head(key string) (*s3.HeadObjectOutput, error) {

	key = conn.prepareKey(key)

	params := &s3.HeadObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
	}

	rsp, err := conn.service.HeadObject(params)

	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (conn *S3Connection) Get(key string) ([]byte, error) {

	key = conn.prepareKey(key)

	params := &s3.GetObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
	}

	rsp, err := conn.service.GetObject(params)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rsp.Body)

	return buf.Bytes(), nil
}

func (conn *S3Connection) Put(key string, body []byte) error {

	key = conn.prepareKey(key)

	params := &s3.PutObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
		ACL:    aws.String("public-read"),
	}

	_, err := conn.service.PutObject(params)

	if err != nil {
		return err
	}

	return nil
}

func (conn *S3Connection) Delete(key string) error {

	key = conn.prepareKey(key)

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
	}

	_, err := conn.service.DeleteObject(params)

	if err != nil {
		return err
	}

	return nil
}

func (conn *S3Connection) prepareKey(key string) string {

	if conn.prefix == "" {
		return key
	}

	return filepath.Join(conn.prefix, key)
}
