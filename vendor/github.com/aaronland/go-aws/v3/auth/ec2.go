package auth

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

func EC2RoleCredentials(ctx context.Context) (aws.Credentials, error) {

	provider := ec2rolecreds.New(func(o *ec2rolecreds.Options) {
		o.Client = imds.New(imds.Options{
			/* custom options */
		})
	})

	return provider.Retrieve(ctx)
}
