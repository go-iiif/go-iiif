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

type S3Thing struct {
	service *s3.S3
	bucket  string
	prefix  string
}

type S3Config struct {
     Bucket string
     Prefix string
     Region string
     Credentials string
}

func NewS3Thing(s3cfg S3Config) (*S3Thing, error) {

	// https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

	sess := session.New(&aws.Config{
		Region:      aws.String(s3cfg.Region),
		Credentials: credentials.NewSharedCredentials("", s3cfg.Credentials),
	})

	_, err := sess.Config.Credentials.Get()

	if err != nil {
		return nil, err
	}

	service := s3.New(sess)

	c := S3Thing{
		service: service,
		bucket:  s3cfg.Bucket,
		prefix:  s3cfg.Prefix,
	}

	return &c, nil
}

func (c *S3Thing) Head(key string) (*s3.HeadObjectOutput, error){

	key = c.prepareKey(key)

	params := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	rsp, err := c.service.HeadObject(params)

	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *S3Thing) Get(key string) ([]byte, error) {

	key = c.prepareKey(key)

	params := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	rsp, err := c.service.GetObject(params)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rsp.Body)

	return buf.Bytes(), nil
}

func (c *S3Thing) Put(key string, body []byte) error {

	key = c.prepareKey(key)

	params := &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
		ACL:    aws.String("public-read"),
	}

	_, err := c.service.PutObject(params)

	if err != nil {
		return err
	}

	return nil
}

func (c *S3Thing) Delete(key string) error {

	key = c.prepareKey(key)

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	_, err := c.service.DeleteObject(params)

	if err != nil {
		return err
	}

	return nil
}

func (c *S3Thing) prepareKey(key string) string {

	if c.prefix == "" {
		return key
	}

	return filepath.Join(c.prefix, key)
}