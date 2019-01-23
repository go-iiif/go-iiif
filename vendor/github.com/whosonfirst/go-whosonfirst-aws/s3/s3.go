package s3

// ideally we could make this conform to some standard "storage" interface
// but that works hasn't been done (20180120/thisisaaronland)

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aaronland/go-string/dsn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
	"github.com/whosonfirst/go-whosonfirst-aws/util"
	"github.com/whosonfirst/go-whosonfirst-mimetypes"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

func ReadCloserFromBytes(b []byte) (io.ReadCloser, error) {
	body := bytes.NewReader(b)
	return ioutil.NopCloser(body), nil
}

type S3Connection struct {
	session  *aws_session.Session
	service  *s3.S3
	uploader *s3manager.Uploader
	bucket   string
	prefix   string
}

type S3Config struct {
	Bucket      string
	Prefix      string
	Region      string
	Credentials string // see notes below
}

type S3ListOptions struct {
	Strict  bool
	Timings bool
	MaxKeys int64
	Path    string
	// Logger log.Logger
}

// this is a nearly straight clone of the core S3 object and
// exists so that people don't have to (re) load all of the
// aws-sdk-go code in their packages (20180801/thisisaaronland)

type S3Object struct {
	KeyRaw       string // what aws-sdk-go returns
	Key          string // KeyRaw but with S3Connection.prefix removed
	Size         int64
	LastModified time.Time
	ETag         string
}

type S3ListCallback func(*S3Object) error

func DefaultS3ListOptions() *S3ListOptions {

	opts := S3ListOptions{
		Strict:  false,
		Timings: false,
		MaxKeys: 500,
	}

	return &opts
}

func ValidS3Credentials() []string {

	valid := []string{
		"env:",
		"iam:",
		"{PROFILE}",
		"{PATH}:{PROFILE}",
	}

	return valid
}

func ValidS3CredentialsString() string {

	valid := ValidS3Credentials()
	return fmt.Sprintf("Valid credential flags are: %s", strings.Join(valid, ", "))
}

func NewS3ConfigFromString(str_dsn string) (*S3Config, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(str_dsn, "bucket", "region", "credentials")

	if err != nil {
		return nil, err
	}

	bucket, _ := dsn_map["bucket"]
	region, _ := dsn_map["region"]
	credentials, _ := dsn_map["credentials"]

	config := S3Config{
		Bucket:      bucket,
		Region:      region,
		Credentials: credentials,
	}

	prefix, ok := dsn_map["prefix"]

	if ok {
		config.Prefix = prefix
	}

	return &config, nil
}

func NewS3Connection(s3cfg *S3Config) (*S3Connection, error) {

	if s3cfg.Bucket == "" {
		return nil, errors.New("Invalid S3 bucket name")
	}

	// https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

	sess, err := session.NewSessionWithCredentials(s3cfg.Credentials, s3cfg.Region)

	if err != nil {
		return nil, err
	}

	service := s3.New(sess)

	uploader := s3manager.NewUploader(sess)

	c := S3Connection{
		session:  sess,
		service:  service,
		uploader: uploader,
		bucket:   s3cfg.Bucket,
		prefix:   s3cfg.Prefix,
	}

	return &c, nil
}

func (conn *S3Connection) URI(key string) string {

	key = conn.prepareKey(key)

	if conn.prefix != "" {
		key = fmt.Sprintf("%s/%s", conn.prefix, key)
	}

	return fmt.Sprintf("https://s3.amazonaws.com/%s/%s", conn.bucket, key)
}

// https://tools.ietf.org/html/rfc7231#section-4.3.2
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectHEAD.html

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

func (conn *S3Connection) Get(key string) (io.ReadCloser, error) {

	key = conn.prepareKey(key)

	params := &s3.GetObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
	}

	rsp, err := conn.service.GetObject(params)

	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}

func (conn *S3Connection) GetBytes(key string) ([]byte, error) {

	fh, err := conn.Get(key)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(fh)

	return buf.Bytes(), nil
}

func (conn *S3Connection) Put(key string, fh io.ReadCloser, args ...interface{}) error {

	// file under known knowns: AWS expects a ReadSeeker for performance
	// and memory reasons but we're passing around ReadClosers  - see also:
	// https://github.com/whosonfirst/go-whosonfirst-readwrite/issues/2
	// (20180120/thisisaaronland)

	defer fh.Close()

	parsed := strings.Split(key, "#")

	key = parsed[0]
	key = conn.prepareKey(key)

	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#UploadInput

	params := s3manager.UploadInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
		Body:   fh,
	}

	ext := filepath.Ext(key)
	types := mimetypes.TypesByExtension(ext)

	if len(types) == 1 {
		params.ContentType = aws.String(types[0])
	}

	// I don't love this... still working it out
	// (20180120/thisisaaronland)

	if len(parsed) > 1 {

		extras := strings.Split(parsed[1], ",")

		for _, ex := range extras {

			kv := strings.Split(ex, "=")

			if len(kv) != 2 {
				return errors.New("Invalid extras")
			}

			k := kv[0]
			v := kv[1]

			switch strings.ToLower(k) {
			case "acl":
				params.ACL = aws.String(v)
			case "contenttype":
				params.ContentType = aws.String(v)
			default:
				// pass
			}
		}
	}

	_, err := conn.uploader.Upload(&params)
	return err
}

func (conn *S3Connection) PutBytes(key string, body []byte) error {

	fh, err := ReadCloserFromBytes(body)

	if err != nil {
		return err
	}

	return conn.Put(key, fh)
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

func (conn *S3Connection) DeleteRecursive(path string) error {

	opts := DefaultS3ListOptions()
	// opts.Timings = *timings
	opts.Path = path

	cb := func(obj *S3Object) error {

		if obj.Key == path {
			return nil
		}

		return conn.DeleteRecursive(obj.Key)
	}

	err := conn.List(cb, opts)

	if err != nil {
		return err
	}

	return conn.Delete(path)
}

func (conn *S3Connection) SetACLForBucket(acl string, opts *S3ListOptions) error {

	cb := func(obj *S3Object) error {

		err := conn.SetACLForKey(obj.Key, acl)
		return err
	}

	return conn.List(cb, opts)
}

func (conn *S3Connection) SetACLForKey(key string, acl string) error {

	key = conn.prepareKey(key)

	params := &s3.PutObjectAclInput{
		ACL:    aws.String(acl),
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
	}

	_, err := conn.service.PutObjectAcl(params)
	return err
}

func (conn *S3Connection) List(cb S3ListCallback, opts *S3ListOptions) error {

	count_pages := int64(0)
	count_items := int64(0)

	if opts.Timings {

		done_ch := make(chan bool)

		defer func() {
			done_ch <- true
		}()

		ticker := time.NewTicker(time.Second * 10)

		go func() {

			for range ticker.C {

				select {
				case <-done_ch:
					break
				default:
					// pass
				}

				log.Printf("items %d pages %d\n", atomic.LoadInt64(&count_items), atomic.LoadInt64(&count_pages))
			}

		}()

		t1 := time.Now()

		defer func() {
			log.Printf("time to list items %d %v\n", atomic.LoadInt64(&count_items), time.Since(t1))
		}()
	}

	prefix := conn.prefix

	if opts.Path != "" {
		prefix = filepath.Join(prefix, opts.Path)
	}

	params := &s3.ListObjectsInput{
		Bucket:  aws.String(conn.bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(opts.MaxKeys),
		// Delimiter: "baz",
	}

	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#ListObjectsOutput
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#Object

	aws_cb := func(rsp *s3.ListObjectsOutput, last_page bool) bool {

		atomic.AddInt64(&count_pages, 1)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done_ch := make(chan bool)
		err_ch := make(chan error)

		for _, aws_obj := range rsp.Contents {

			atomic.AddInt64(&count_items, 1)

			// because this:
			// https://github.com/whosonfirst/go-whosonfirst-aws/issues/1

			key_raw := *aws_obj.Key
			key := key_raw

			if conn.prefix != "" {

				prefix := fmt.Sprintf("%s/", conn.prefix)

				if strings.HasPrefix(key, prefix) {
					key = strings.Replace(key, prefix, "", -1)
				}
			}

			obj := &S3Object{
				KeyRaw:       key_raw,
				Key:          key,
				Size:         *aws_obj.Size,
				ETag:         *aws_obj.ETag,
				LastModified: *aws_obj.LastModified,
			}

			go func(ctx context.Context, obj *S3Object, done_ch chan bool, err_ch chan error) {

				defer func() {
					done_ch <- true
				}()

				select {
				case <-ctx.Done():
					return
				default:
					// pass
				}

				err := cb(obj)

				if err != nil {
					msg := fmt.Sprintf("failed to process %s because %s", obj.Key, err)
					err_ch <- errors.New(msg)
				}

			}(ctx, obj, done_ch, err_ch)
		}

		remaining := len(rsp.Contents)
		ok := true

		for remaining > 0 {

			select {

			case <-done_ch:
				remaining -= 1
			case e := <-err_ch:
				log.Println(e)

				/*
					if opts.Strict {
						ok = false
						break
					}
				*/

			default:
				// pass
			}
		}

		return ok
	}

	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#example_S3_ListObjects_shared00

	err := conn.service.ListObjectsPages(params, aws_cb)

	if err != nil {
		return err
	}

	return nil
}

func (conn *S3Connection) HasChanged(key string, local []byte) (bool, error) {

	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#HeadObjectInput
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#HeadObjectOutput

	head, err := conn.Head(key)

	if err != nil {

		aws_err := err.(awserr.Error)

		if aws_err.Code() == "NotFound" {
			return true, nil
		}

		if aws_err.Code() == "SlowDown" {

		}

		return false, err
	}

	enc := md5.Sum(local)
	local_hash := hex.EncodeToString(enc[:])

	etag := *head.ETag
	remote_hash := strings.Replace(etag, "\"", "", -1)

	if local_hash == remote_hash {
		return false, nil
	}

	return true, nil
}

func IsNotFound(err error) bool {
	return util.IsAWSErrorWithCode(err, s3.ErrCodeNoSuchKey)
}

func (conn *S3Connection) prepareKey(key string) string {

	if conn.prefix == "" {
		return key
	}

	return filepath.Join(conn.prefix, key)
}
