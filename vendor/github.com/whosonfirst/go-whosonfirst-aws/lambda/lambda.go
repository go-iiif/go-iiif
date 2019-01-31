package lambda

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	aws_lambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
)

func NewLambdaServiceWithDSN(str_dsn string) (*aws_lambda.Lambda, error) {

	sess, err := session.NewSessionWithDSN(str_dsn)

	if err != nil {
		return nil, err
	}

	svc := aws_lambda.New(sess)
	return svc, nil
}

func InvokeFunction(svc *aws_lambda.Lambda, lambda_func string, lambda_type string, payload interface{}) (*aws_lambda.InvokeOutput, error) {

	enc_payload, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	input := &aws_lambda.InvokeInput{
		FunctionName:   aws.String(lambda_func),
		InvocationType: aws.String(lambda_type),
		Payload:        enc_payload,
	}

	if *input.InvocationType == "RequestResponse" {
		input.LogType = aws.String("Tail")
	}

	rsp, err := svc.Invoke(input)

	if err != nil {
		return nil, err
	}

	if *input.InvocationType != "RequestResponse" {
		return nil, nil
	}

	enc_result := *rsp.LogResult

	result, err := base64.StdEncoding.DecodeString(enc_result)

	if err != nil {
		return nil, err
	}

	if *rsp.StatusCode != 200 {
		return nil, errors.New(string(result))
	}

	return rsp, nil
}
