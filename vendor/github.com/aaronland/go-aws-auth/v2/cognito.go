package auth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

// STSCredentialsForDeveloperIdentityOptions defines configuration options for generating temporary STS (AWS) credentials for a developer identity.
type STSCredentialsForDeveloperIdentityOptions struct {
	// An identity pool ID in the format REGION:GUID.
	IdentityPoolId string
	// A set of name-value pairs that map developer or provider names to a user.
	Logins map[string]string
	// The Amazon Resource Name (ARN) of the role that the caller is assuming.
	RoleArn string
	// An identifier for the assumed role session.
	RoleSessionName string
	// The duration, in seconds, of the role session.
	Duration int32
	// An optional list of Amazon Resource Names (ARNs)  that you want to use as managed session policies.
	Policies []string
}

// STSCredentialsForDeveloperIdentity generate temporary STS (AWS) credentials for a developer identity.
func STSCredentialsForDeveloperIdentity(ctx context.Context, aws_cfg aws.Config, opts *STSCredentialsForDeveloperIdentityOptions) (*types.Credentials, error) {

	cognito_client := cognitoidentity.NewFromConfig(aws_cfg)
	sts_client := sts.NewFromConfig(aws_cfg)

	//

	// Get temporary OpenID token from Cognito

	token_opts := &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
		IdentityPoolId: aws.String(opts.IdentityPoolId),
		Logins:         opts.Logins,
	}

	token_rsp, err := cognito_client.GetOpenIdTokenForDeveloperIdentity(ctx, token_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive token for developer identity, %w", err)
	}

	// Get temporary credentials from STS

	creds_opts := &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(opts.RoleArn),
		RoleSessionName:  aws.String(opts.RoleSessionName),
		WebIdentityToken: token_rsp.Token,
		DurationSeconds:  aws.Int32(opts.Duration),
	}

	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sts#AssumeRoleWithWebIdentityInput
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html#policies_session

	if len(opts.Policies) > 0 {

		session_policies := make([]types.PolicyDescriptorType, len(opts.Policies))

		for idx, arn := range opts.Policies {

			session_policies[idx] = types.PolicyDescriptorType{
				Arn: aws.String(arn),
			}
		}

		creds_opts.PolicyArns = session_policies
	}

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_request.html

	creds_rsp, err := sts_client.AssumeRoleWithWebIdentity(ctx, creds_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to assume role, %w", err)
	}

	return creds_rsp.Credentials, nil
}
