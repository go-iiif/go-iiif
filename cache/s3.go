package cache

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/thisisaaronland/go-iiif/config"
	_ "log"
	"path/filepath"
	"strings"
)

type S3Cache struct {
	service *s3.S3
	bucket  string
	prefix  string
}

func NewS3Cache(cfg config.CacheConfig) (*S3Cache, error) {

	// https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html

	bucket := cfg.Path
	prefix := ""

	parts := strings.Split(bucket, "/")

	if len(parts) > 1 {
		bucket = parts[0]
		prefix = strings.Join(parts[1:], "/")
	}

	// please read from the config (20160923/thisisaaronland)

	region := "us-east-1"
	creds := "default"

	sess := session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", creds),
	})

	_, err := sess.Config.Credentials.Get()

	if err != nil {
		return nil, err
	}

	service := s3.New(sess)

	c := S3Cache{
		service: service,
		bucket:  bucket,
		prefix:  prefix,
	}

	return &c, nil
}

func (c *S3Cache) Exists(key string) bool {

	key = c.prepareKey(key)

	params := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	_, err := c.service.HeadObject(params)

	if err != nil {
		return false
	}

	return true
}

func (c *S3Cache) Get(key string) ([]byte, error) {

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

func (c *S3Cache) Set(key string, body []byte) error {

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

func (c *S3Cache) Unset(key string) error {

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

func (c *S3Cache) prepareKey(key string) string {

	if c.prefix == "" {
		return key
	}

	return filepath.Join(c.prefix, key)
}