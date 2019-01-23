package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-aws/s3"
	"io"
	_ "io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	valid_flags := s3.ValidS3CredentialsString()

	var dsn = flag.String("dsn", "", "A valid S3 DSN")

	var bucket = flag.String("bucket", "", "A valid S3 bucket name")
	var prefix = flag.String("prefix", "", "An optional path/prefix inside the S3 bucket")
	var region = flag.String("region", "us-east-1", "A valid AWS S3 region")
	var credentials = flag.String("credentials", "env:", "A valid S3 credentials flag. "+valid_flags)

	var recursive = flag.Bool("recursive", false, "...")
	var timings = flag.Bool("timings", false, "...")

	flag.Parse()

	var config *s3.S3Config

	if *dsn != "" {

		c, err := s3.NewS3ConfigFromString(*dsn)

		if err != nil {
			log.Fatal(err)
		}

		config = c
	} else {

		config = &s3.S3Config{
			Bucket:      *bucket,
			Prefix:      *prefix,
			Region:      *region,
			Credentials: *credentials,
		}
	}

	conn, err := s3.NewS3Connection(config)

	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()
	count := len(args)

	if count == 0 {
		log.Fatal("Missing S3 command")
	}

	if count == 1 {
		log.Fatal("Missing anything to do with your S3 command")
	}

	cmd := strings.ToUpper(args[0])

	for _, path := range args[1:] {

		var rsp interface{}
		var err error

		switch cmd {

		case "ACL":

			parsed := strings.Split(path, "#")

			if len(parsed) != 2 {
				log.Fatal("Invalid acl string")
			}

			key := parsed[0]
			acl := parsed[1]

			if acl == "" {
				log.Fatal("Missing ACL")
			}

			if key == "" {

				opts := s3.DefaultS3ListOptions()
				opts.Timings = *timings

				err = conn.SetACLForBucket(acl, opts)
			} else {
				err = conn.SetACLForKey(key, acl)
			}

			log.Println(parsed, len(parsed))

		case "HEAD":

			rsp, err = conn.Head(path)

		case "GET":

			rsp, err = conn.Get(path)

		case "LIST":

			opts := s3.DefaultS3ListOptions()
			opts.Timings = *timings

			cb := func(obj *s3.S3Object) error {

				log.Printf("%s (%s)\n", obj.Key, obj.KeyRaw)
				return nil
			}

			err = conn.List(cb, opts)

		case "PUT":

			parsed := strings.Split(path, "#")

			rel_path := parsed[0]
			abs_path, err := filepath.Abs(rel_path)

			if err != nil {
				log.Fatal(err)
			}

			key := filepath.Base(abs_path)

			if len(parsed) == 2 {
				key = fmt.Sprintf("%s#%s", key, parsed[1])
			}

			fh, err := os.Open(abs_path)

			if err != nil {
				log.Fatal(err)
			}

			err = conn.Put(key, fh)

		case "DELETE":

			if *recursive {
				err = conn.DeleteRecursive(path)

			} else {
				err = conn.Delete(path)
			}

		case "URI":

			rsp = conn.URI(path)

		default:

			log.Fatal("Invalid command")
		}

		if err != nil {
			log.Fatal(err)
		}

		switch cmd {

		case "HEAD":
			log.Println(rsp)

		case "URI":

			uri := rsp.(string)
			os.Stdout.Write([]byte(uri))

		case "GET":

			fh := rsp.(io.ReadCloser)

			_, err := io.Copy(os.Stdout, fh)

			if err != nil {
				log.Fatal(err)
			}

		default:
			// pass
		}
	}

	os.Exit(0)
}
