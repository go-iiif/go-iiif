package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/url"

	aa_bucket "github.com/aaronland/gocloud-blob/bucket"
	aa_s3 "github.com/aaronland/gocloud-blob/s3"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {

	var source_bucket_uri string
	var target_bucket_uri string
	var source_key string
	var target_key string
	var str_acl string

	flag.StringVar(&source_bucket_uri, "source-bucket-uri", "", "A valid gocloud.dev/blob.Bucket URI.")
	flag.StringVar(&target_bucket_uri, "target-bucket-uri", "", "A valid gocloud.dev/blob.Bucket URI.")
	flag.StringVar(&source_key, "source-key", "", "The key to read from the source bucket")
	flag.StringVar(&target_key, "target-key", "", "The key to write to the target bucket")
	flag.StringVar(&str_acl, "acl", "public-read", "A valid AWS S3 ACL string.")

	flag.Parse()

	ctx := context.Background()

	slog.Info("Create source bucket", "uri", source_bucket_uri)

	source_bucket, err := aa_bucket.OpenBucket(ctx, source_bucket_uri)
	// source_bucket, err := blob.OpenBucket(ctx, source_bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open source bucket, %v", err)
	}

	defer source_bucket.Close()

	slog.Info("Create target bucket", "uri", target_bucket_uri)

	target_bucket, err := aa_bucket.OpenBucket(ctx, target_bucket_uri)
	// target_bucket, err := blob.OpenBucket(ctx, target_bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open target bucket, %v", err)
	}

	defer target_bucket.Close()

	r, err := source_bucket.NewReader(ctx, source_key, nil)

	if err != nil {
		log.Fatalf("Failed to create new reader for source key, %v", err)
	}

	defer r.Close()

	acl, err := aa_s3.StringACLToObjectCannedACL(str_acl)

	if err != nil {
		log.Fatalf("Failed to derive ACL object, %v", err)
	}

	before := func(asFunc func(interface{}) bool) error {

		req := &aws_s3.PutObjectInput{}
		ok := asFunc(&req)

		if !ok {
			return fmt.Errorf("Not a PutObjectInput instance")
		}

		req.ACL = acl

		b, _ := url.Parse(target_bucket_uri)

		slog.Info("Before write bucket test", "target", b.Host, "put object", *req.Bucket, "key", *req.Key)

		if *req.Bucket != b.Host {
			return fmt.Errorf("Request bucket does not match target bucket")
		}

		return nil
	}

	wr_opts := &blob.WriterOptions{
		BeforeWrite: before,
	}

	slog.Info("Create target writer", "key", target_key)

	wr, err := target_bucket.NewWriter(ctx, target_key, wr_opts)

	if err != nil {
		log.Fatalf("Failed to create new writer for target key, %v", err)
	}

	_, err = io.Copy(wr, r)

	if err != nil {
		log.Fatalf("Failed to copy source key to target key, %v", err)
	}

	err = wr.Close()

	if err != nil {
		log.Fatalf("Failed to close target writer, %v", err)
	}
}
