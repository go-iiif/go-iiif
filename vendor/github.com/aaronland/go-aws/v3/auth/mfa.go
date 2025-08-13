package auth

// https://docs.aws.amazon.com/cli/latest/reference/sts/get-session-token.html

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

// GetCredentialsWithMFAWithContext will return a time-limited AWS credentials (defined by 'duration') for the account associated with 'cfg'
// a valid multi-factor authentication (MFS) token defined by 'token' using a custom `context.Context` instance.
func GetCredentialsWithMFAWithContext(ctx context.Context, cfg aws.Config, device string, token string, duration int32) (*types.Credentials, error) {

	sts_client := sts.NewFromConfig(cfg)

	username, err := username(ctx, sts_client)

	if err != nil {
		return nil, err
	}

	iam_client := iam.NewFromConfig(cfg)

	serial_number, err := mfaDevice(ctx, iam_client, username, device)

	if err != nil {
		return nil, err
	}

	return sessionCredentials(ctx, sts_client, serial_number, token, duration)
}

func username(ctx context.Context, sts_client *sts.Client) (string, error) {

	opts := &sts.GetCallerIdentityInput{}

	rsp, err := sts_client.GetCallerIdentity(ctx, opts)

	if err != nil {
		return "", fmt.Errorf("Failed to derive username, %w", err)
	}

	arn := rsp.Arn

	if strings.Contains(*arn, ":user/") {
		return strings.Split(*arn, ":user/")[1], nil
	}

	/*
		if strings.Contains(*arn, ":assumed-role/") {
			return strings.Split(*arn, ":assumed-role/")[1], nil
		}
	*/

	return "", fmt.Errorf("Failed to derive user name from ARN")
}

func mfaDevice(ctx context.Context, iam_client *iam.Client, user_arn string, serial_number string) (string, error) {

	opts := &iam.ListMFADevicesInput{
		UserName: &user_arn,
	}

	rsp, err := iam_client.ListMFADevices(ctx, opts)

	if err != nil {
		return "", fmt.Errorf("Failed to list devices for %s, %w", user_arn, err)
	}

	count_devices := len(rsp.MFADevices)

	if count_devices == 0 {
		return "", fmt.Errorf("No devices found")
	}

	if serial_number == "" {

		if count_devices > 1 {
			return "", fmt.Errorf("Multiple devices found, need serial number")
		}

		return *rsp.MFADevices[0].SerialNumber, nil
	}

	for _, d := range rsp.MFADevices {

		if *d.SerialNumber == serial_number {
			return serial_number, nil
		}
	}

	return "", fmt.Errorf("No matching serial number")
}

func sessionCredentials(ctx context.Context, sts_client *sts.Client, mfaDevice string, tokenCode string, duration int32) (*types.Credentials, error) {

	opts := &sts.GetSessionTokenInput{
		SerialNumber:    &mfaDevice,
		DurationSeconds: &duration,
		TokenCode:       &tokenCode,
	}

	rsp, err := sts_client.GetSessionToken(ctx, opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to get session token, %w", err)
	}

	creds := rsp.Credentials
	return creds, nil
}
