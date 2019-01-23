package sqs

import (
	"github.com/aaronland/go-string/dsn"
	"github.com/aws/aws-sdk-go/aws"
	aws_sqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
	"strings"
)

func NewSQSServiceWithDSN(str_dsn string) (*aws_sqs.SQS, string, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(str_dsn, "region", "credentials", "queue")

	if err != nil {
		return nil, "", err
	}

	sqs_creds, _ := dsn_map["credentials"]
	sqs_region, _ := dsn_map["region"]
	sqs_queue, _ := dsn_map["queue"]

	sess, err := session.NewSessionWithCredentials(sqs_creds, sqs_region)

	if err != nil {
		return nil, "", err
	}

	svc := aws_sqs.New(sess)

	if !strings.HasPrefix(sqs_queue, "https://sqs") {

		rsp, err := svc.GetQueueUrl(&aws_sqs.GetQueueUrlInput{
			QueueName: aws.String(sqs_queue),
		})

		if err != nil {
			return nil, "", err
		}

		sqs_queue = *rsp.QueueUrl
	}

	return svc, sqs_queue, nil
}
