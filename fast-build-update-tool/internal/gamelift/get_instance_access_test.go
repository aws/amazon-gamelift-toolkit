package gamelift

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	"github.com/aws/aws-sdk-go-v2/service/gamelift/types"
	"github.com/stretchr/testify/assert"
)

var (
	accessKey       = "accessKey"
	secretAccessKey = "secretAccessKey"
	sessionToken    = "sessionToken"
	instanceId      = "i-12345"
	fleetId         = "fleet-6789"
)

// TestGetInstanceAccessSuccess verifies that GetComputeAccess is called with the proper inputs, and returns the proper data
func TestGetInstanceAccessSuccess(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	// Ensure get compute access is called with proper params
	awsMock.GetComputeAccessFunc = func(ctx context.Context, params *gamelift.GetComputeAccessInput, optFns ...func(*gamelift.Options)) (*gamelift.GetComputeAccessOutput, error) {
		return &gamelift.GetComputeAccessOutput{
			Credentials: &types.AwsCredentials{
				AccessKeyId:     &accessKey,
				SecretAccessKey: &secretAccessKey,
				SessionToken:    &sessionToken,
			},
		}, nil
	}

	output, err := client.GetInstanceAccess(context.Background(), fleetId, instanceId)

	assert.Nil(t, err)

	// Ensure we call with the proper params
	assert.Equal(t, fleetId, *awsMock.GetComputeAccessCalls()[0].Params.FleetId)
	assert.Equal(t, instanceId, *awsMock.GetComputeAccessCalls()[0].Params.ComputeName)

	// Ensure we return the proper output
	assert.Equal(t, accessKey, output.AccessKeyId)
	assert.Equal(t, secretAccessKey, output.SecretAccessKey)
	assert.Equal(t, sessionToken, output.SessionToken)
}

// TestGetInstanceAccessError ensures that any errors called with GetComputeAccess are handled properly
func TestGetInstanceAccessError(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	expectedError := errors.New("test error")

	awsMock.GetComputeAccessFunc = func(ctx context.Context, params *gamelift.GetComputeAccessInput, optFns ...func(*gamelift.Options)) (*gamelift.GetComputeAccessOutput, error) {
		return nil, expectedError
	}

	_, err := client.GetInstanceAccess(context.Background(), fleetId, instanceId)

	// Ensure we bubble up any errors
	assert.Equal(t, expectedError, err)
}
