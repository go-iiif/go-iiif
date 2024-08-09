// Code generated by smithy-go-codegen DO NOT EDIT.

package iam

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Sets the specified version of the global endpoint token as the token version
// used for the Amazon Web Services account.
//
// By default, Security Token Service (STS) is available as a global service, and
// all STS requests go to a single endpoint at https://sts.amazonaws.com . Amazon
// Web Services recommends using Regional STS endpoints to reduce latency, build in
// redundancy, and increase session token availability. For information about
// Regional endpoints for STS, see [Security Token Service endpoints and quotas]in the Amazon Web Services General Reference.
//
// If you make an STS call to the global endpoint, the resulting session tokens
// might be valid in some Regions but not others. It depends on the version that is
// set in this operation. Version 1 tokens are valid only in Amazon Web Services
// Regions that are available by default. These tokens do not work in manually
// enabled Regions, such as Asia Pacific (Hong Kong). Version 2 tokens are valid in
// all Regions. However, version 2 tokens are longer and might affect systems where
// you temporarily store tokens. For information, see [Activating and deactivating STS in an Amazon Web Services Region]in the IAM User Guide.
//
// To view the current session token version, see the GlobalEndpointTokenVersion
// entry in the response of the GetAccountSummaryoperation.
//
// [Activating and deactivating STS in an Amazon Web Services Region]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html
// [Security Token Service endpoints and quotas]: https://docs.aws.amazon.com/general/latest/gr/sts.html
func (c *Client) SetSecurityTokenServicePreferences(ctx context.Context, params *SetSecurityTokenServicePreferencesInput, optFns ...func(*Options)) (*SetSecurityTokenServicePreferencesOutput, error) {
	if params == nil {
		params = &SetSecurityTokenServicePreferencesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "SetSecurityTokenServicePreferences", params, optFns, c.addOperationSetSecurityTokenServicePreferencesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*SetSecurityTokenServicePreferencesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type SetSecurityTokenServicePreferencesInput struct {

	// The version of the global endpoint token. Version 1 tokens are valid only in
	// Amazon Web Services Regions that are available by default. These tokens do not
	// work in manually enabled Regions, such as Asia Pacific (Hong Kong). Version 2
	// tokens are valid in all Regions. However, version 2 tokens are longer and might
	// affect systems where you temporarily store tokens.
	//
	// For information, see [Activating and deactivating STS in an Amazon Web Services Region] in the IAM User Guide.
	//
	// [Activating and deactivating STS in an Amazon Web Services Region]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html
	//
	// This member is required.
	GlobalEndpointTokenVersion types.GlobalEndpointTokenVersion

	noSmithyDocumentSerde
}

type SetSecurityTokenServicePreferencesOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationSetSecurityTokenServicePreferencesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpSetSecurityTokenServicePreferences{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpSetSecurityTokenServicePreferences{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "SetSecurityTokenServicePreferences"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpSetSecurityTokenServicePreferencesValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opSetSecurityTokenServicePreferences(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opSetSecurityTokenServicePreferences(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "SetSecurityTokenServicePreferences",
	}
}
