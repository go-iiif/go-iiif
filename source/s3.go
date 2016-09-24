package source

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
)

type S3Source struct {
	service *s3.S3
	bucket  string
}

func NewS3Source(cfg *iiifconfig.Config) (*S3Source, error) {

	// grab stuff from cache/s3.go (20160923/thisisaaronland)

	session, err := session.NewSession()

	if err != nil {
		return nil, err
	}

	service := s3.New(session)

	bucket := "fixme"

	c := S3Source{
		service: service,
		bucket:  bucket,
	}

	return &c, nil
}

func (c *S3Source) Read(id string) ([]byte, error) {

	params := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(id),
	}

	rsp, err := c.service.GetObject(params)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rsp.Body)

	return buf.Bytes(), nil
}
