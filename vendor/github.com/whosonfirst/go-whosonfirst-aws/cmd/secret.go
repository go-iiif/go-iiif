package main

// https://docs.aws.amazon.com/sdk-for-go/api/service/secretsmanager/#example_SecretsManager_GetSecretValue_shared00

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/whosonfirst/go-whosonfirst-aws/config"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
	"log"
)

func main() {

	valid_creds := config.ValidCredentialsString()

	var region = flag.String("region", "us-east-1", "A valid AWS region")
	var credentials = flag.String("credentials", "env:", "A valid AWS credentials flag. "+valid_creds)

	flag.Parse()

	sess, err := session.NewSessionWithCredentials(*credentials, *region)

	if err != nil {
		log.Fatal(err)
	}

	// please for to wrap me in a simplified interface/function...
	// (20180925/thisisaaronland)

	svc := secretsmanager.New(sess)

	for _, secret_name := range flag.Args() {

		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secret_name),
		}

		result, err := svc.GetSecretValue(input)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	}

}
