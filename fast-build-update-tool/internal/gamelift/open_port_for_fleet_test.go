package gamelift

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	"github.com/aws/aws-sdk-go-v2/service/gamelift/types"
	"github.com/stretchr/testify/assert"
)

// TestOpenPortForFleet ensures we call UpdateFleetPortSettingsFunc with the proper params
func TestOpenPortForFleet(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	port := int32(1026)
	ipRange := "0.0.0.0/0"

	awsMock.UpdateFleetPortSettingsFunc = func(ctx context.Context, params *gamelift.UpdateFleetPortSettingsInput, optFns ...func(*gamelift.Options)) (*gamelift.UpdateFleetPortSettingsOutput, error) {
		return nil, nil
	}

	err := client.OpenPortForFleet(context.Background(), fleetId, port, ipRange)

	assert.Nil(t, err)

	calls := awsMock.UpdateFleetPortSettingsCalls()
	assert.Len(t, calls, 1)
	assert.Equal(t, fleetId, *calls[0].Params.FleetId)
	assert.Equal(t, ipRange, *calls[0].Params.InboundPermissionAuthorizations[0].IpRange)
	assert.Equal(t, port, *calls[0].Params.InboundPermissionAuthorizations[0].FromPort)
	assert.Equal(t, port, *calls[0].Params.InboundPermissionAuthorizations[0].ToPort)
}

// TestOpenPortForFleetHandleDuplicateError ensures that we don't error when a port is already open for a fleet
func TestOpenPortForFleetHandleDuplicateError(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	awsMock.UpdateFleetPortSettingsFunc = func(ctx context.Context, params *gamelift.UpdateFleetPortSettingsInput, optFns ...func(*gamelift.Options)) (*gamelift.UpdateFleetPortSettingsOutput, error) {
		return nil, &types.InvalidRequestException{Message: aws.String("InvalidPermission.Duplicate error")}
	}

	err := client.OpenPortForFleet(context.Background(), fleetId, 1026, "127.0.0.1/32")

	// Ensure we don't forward this error
	assert.Nil(t, err)
	assert.Len(t, awsMock.UpdateFleetPortSettingsCalls(), 1)
}
