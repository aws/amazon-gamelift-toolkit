// gamelift contains any logic around interacting with the AWS GameLift service through the AWS SDK
package gamelift

import (
	"context"

	"github.com/aws/smithy-go/logging"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
)

// GameLiftClient is uses to manage any direct interactions with the AWS GameLift service
type GameLiftClient struct {
	gamelift AWSGameliftClient
}

// NewGameLiftClient will build a new GameLiftClient with the default AWS credentials
func NewGameLiftClient(ctx context.Context, logger logging.Logger) (*GameLiftClient, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return &GameLiftClient{gamelift: gamelift.NewFromConfig(cfg)}, nil
}

//go:generate moq -skip-ensure -out ./moq_aws_gamelift_client_test.go  . AWSGameliftClient

// AWSGameliftClient wraps the expected GameLift interface from the AWS SDK
type AWSGameliftClient interface {
	DescribeFleetAttributes(ctx context.Context, params *gamelift.DescribeFleetAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetAttributesOutput, error)
	DescribeRuntimeConfiguration(ctx context.Context, params *gamelift.DescribeRuntimeConfigurationInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeRuntimeConfigurationOutput, error)
	UpdateFleetPortSettings(ctx context.Context, params *gamelift.UpdateFleetPortSettingsInput, optFns ...func(*gamelift.Options)) (*gamelift.UpdateFleetPortSettingsOutput, error)
	DescribeFleetLocationAttributes(ctx context.Context, params *gamelift.DescribeFleetLocationAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetLocationAttributesOutput, error)
	DescribeInstances(ctx context.Context, params *gamelift.DescribeInstancesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeInstancesOutput, error)
	GetComputeAccess(ctx context.Context, params *gamelift.GetComputeAccessInput, optFns ...func(*gamelift.Options)) (*gamelift.GetComputeAccessOutput, error)
}
