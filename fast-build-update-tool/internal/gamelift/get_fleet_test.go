package gamelift

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	"github.com/aws/aws-sdk-go-v2/service/gamelift/types"
	"github.com/stretchr/testify/assert"
)

// TestGetFleet tests the golden path of fleet lookup functionality
func TestGetFleet(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	expectedExeOne := "/local/game/myserver.exe"
	expectedExeTwo := "/local/game/launcher.exe"

	awsMock.DescribeFleetAttributesFunc = func(ctx context.Context, params *gamelift.DescribeFleetAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetAttributesOutput, error) {
		return &gamelift.DescribeFleetAttributesOutput{
			FleetAttributes: []types.FleetAttributes{
				types.FleetAttributes{OperatingSystem: types.OperatingSystemAmazonLinux},
			},
		}, nil
	}

	awsMock.DescribeRuntimeConfigurationFunc = func(ctx context.Context, params *gamelift.DescribeRuntimeConfigurationInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeRuntimeConfigurationOutput, error) {
		return &gamelift.DescribeRuntimeConfigurationOutput{
			RuntimeConfiguration: &types.RuntimeConfiguration{
				ServerProcesses: []types.ServerProcess{
					types.ServerProcess{LaunchPath: aws.String(expectedExeOne)},
					types.ServerProcess{LaunchPath: aws.String(expectedExeOne)}, // Make sure we filter out any dupes
					types.ServerProcess{LaunchPath: aws.String(expectedExeTwo)},
				},
			},
		}, nil
	}

	fleet, err := client.GetFleet(context.Background(), fleetId)
	assert.Nil(t, err)

	describeFleetCalls := awsMock.DescribeFleetAttributesCalls()
	assert.Len(t, describeFleetCalls, 1)
	assert.Len(t, describeFleetCalls[0].Params.FleetIds, 1)
	assert.Equal(t, fleetId, describeFleetCalls[0].Params.FleetIds[0])

	describeRuntimeCalls := awsMock.DescribeRuntimeConfigurationCalls()
	assert.Len(t, describeRuntimeCalls, 1)
	assert.Equal(t, fleetId, *describeRuntimeCalls[0].Params.FleetId)

	assert.Equal(t, fleetId, fleet.Id)
	assert.Equal(t, config.OperatingSystemLinux, fleet.OperatingSystem)
	assert.Len(t, fleet.ExecutablePaths, 2)
	assert.Contains(t, fleet.ExecutablePaths, expectedExeOne)
	assert.Contains(t, fleet.ExecutablePaths, expectedExeTwo)
}

// TestGetFleetLookupError tests the case where DescribeFleetAttributes returns an error
func TestGetFleetLookupError(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	expected := errors.New("test error")

	awsMock.DescribeFleetAttributesFunc = func(ctx context.Context, params *gamelift.DescribeFleetAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetAttributesOutput, error) {
		return nil, expected
	}

	_, err := client.GetFleet(context.Background(), fleetId)
	assert.ErrorContains(t, err, expected.Error())
}

// TestGetFleetLookupErrorEmptyAttributes tests the case where DescribeFleetAttributes returns an empty slice
func TestGetFleetLookupErrorEmptyAttributes(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	awsMock.DescribeFleetAttributesFunc = func(ctx context.Context, params *gamelift.DescribeFleetAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetAttributesOutput, error) {
		return &gamelift.DescribeFleetAttributesOutput{FleetAttributes: []types.FleetAttributes{}}, nil
	}

	_, err := client.GetFleet(context.Background(), fleetId)
	assert.ErrorContains(t, err, "fleet not found")
}

// TestGetFleetLookupErrorEmptyAttributes tests the case where DescribeRuntimeConfiguration returns an error
func TestGetFleetRuntimeLookupError(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	expected := errors.New("test error")

	awsMock.DescribeFleetAttributesFunc = func(ctx context.Context, params *gamelift.DescribeFleetAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetAttributesOutput, error) {
		return &gamelift.DescribeFleetAttributesOutput{
			FleetAttributes: []types.FleetAttributes{
				types.FleetAttributes{OperatingSystem: types.OperatingSystemAmazonLinux},
			},
		}, nil
	}

	awsMock.DescribeRuntimeConfigurationFunc = func(ctx context.Context, params *gamelift.DescribeRuntimeConfigurationInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeRuntimeConfigurationOutput, error) {
		return nil, expected
	}

	_, err := client.GetFleet(context.Background(), fleetId)
	assert.ErrorContains(t, err, expected.Error())
}
