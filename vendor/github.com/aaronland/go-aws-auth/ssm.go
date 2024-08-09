package auth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func NewSSMClient(ctx context.Context, uri string) (*ssm.Client, error) {

	cfg, err := NewConfig(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new config, %w", err)
	}

	return ssm.NewFromConfig(cfg), nil
}

func NewSSMClientWithCredentialsString(ctx context.Context, str_creds string) (*ssm.Client, error) {

	cfg, err := NewConfigWithCredentialsString(ctx, str_creds)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new config, %w", err)
	}

	return ssm.NewFromConfig(cfg), nil
}
