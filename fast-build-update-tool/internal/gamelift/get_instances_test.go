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

// TestGetInstances verifies that we call the proper AWS functions and filter out inactive results when looking up an instance
func TestGetInstances(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}
	badRegion := "eu-west-1"

	// Set up a fleet with two active regions (will be used), and one activating region (will be ignored)
	awsMock.DescribeFleetLocationAttributesFunc = func(ctx context.Context, params *gamelift.DescribeFleetLocationAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetLocationAttributesOutput, error) {
		return &gamelift.DescribeFleetLocationAttributesOutput{
			LocationAttributes: []types.LocationAttributes{
				types.LocationAttributes{LocationState: &types.LocationState{Location: aws.String("us-west-1"), Status: types.FleetStatusActive}},
				types.LocationAttributes{LocationState: &types.LocationState{Location: aws.String(badRegion), Status: types.FleetStatusActivating}},
				types.LocationAttributes{LocationState: &types.LocationState{Location: aws.String("us-east-1"), Status: types.FleetStatusActive}},
			},
		}, nil
	}

	// Set up describe instances to return three instances per region, two active and one inactive.
	// Inactive will be filtered out, the other will be filtered out by id
	awsMock.DescribeInstancesFunc = func(ctx context.Context, params *gamelift.DescribeInstancesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeInstancesOutput, error) {
		if *params.Location == "us-west-1" {
			return &gamelift.DescribeInstancesOutput{
				Instances: []types.Instance{
					types.Instance{OperatingSystem: types.OperatingSystemAmazonLinux2023, IpAddress: aws.String("127.0.0.1"), InstanceId: aws.String("us-west-instance-1"), Location: params.Location, Status: types.InstanceStatusPending},
					types.Instance{OperatingSystem: types.OperatingSystemAmazonLinux2023, IpAddress: aws.String("127.0.0.2"), InstanceId: aws.String("us-west-instance-2"), Location: params.Location, Status: types.InstanceStatusActive},
					types.Instance{OperatingSystem: types.OperatingSystemAmazonLinux2023, IpAddress: aws.String("127.0.0.2"), InstanceId: aws.String("us-west-instance-3"), Location: params.Location, Status: types.InstanceStatusActive},
				},
			}, nil
		} else if *params.Location == "us-east-1" {
			return &gamelift.DescribeInstancesOutput{
				Instances: []types.Instance{
					types.Instance{OperatingSystem: types.OperatingSystemWindows2016, IpAddress: aws.String("192.0.0.1"), InstanceId: aws.String("us-east-instance-1"), Location: params.Location, Status: types.InstanceStatusPending},
					types.Instance{OperatingSystem: types.OperatingSystemWindows2016, IpAddress: aws.String("192.0.0.2"), InstanceId: aws.String("us-east-instance-2"), Location: params.Location, Status: types.InstanceStatusActive},
					types.Instance{OperatingSystem: types.OperatingSystemWindows2016, IpAddress: aws.String("192.0.0.2"), InstanceId: aws.String("us-east-instance-3"), Location: params.Location, Status: types.InstanceStatusActive},
				},
			}, nil
		}
		return nil, nil
	}

	instances, err := client.GetInstances(context.Background(), fleetId, []string{"us-west-instance-2", "us-east-instance-2"})

	assert.Nil(t, err)
	assert.Len(t, instances, 2)

	assert.Len(t, awsMock.DescribeFleetLocationAttributesCalls(), 1)
	assert.Equal(t, fleetId, *awsMock.DescribeFleetLocationAttributesCalls()[0].Params.FleetId)

	assert.Len(t, awsMock.DescribeInstancesCalls(), 2)
	assert.Equal(t, "us-west-1", *awsMock.DescribeInstancesCalls()[0].Params.Location)
	assert.Equal(t, "us-east-1", *awsMock.DescribeInstancesCalls()[1].Params.Location)

	instanceOne := instances[0]
	assert.Equal(t, "127.0.0.2", instanceOne.IpAddress)
	assert.Equal(t, "us-west-instance-2", instanceOne.InstanceId)
	assert.Equal(t, "us-west-1", instanceOne.Region)
	assert.Equal(t, fleetId, instanceOne.FleetId)
	assert.Equal(t, config.OperatingSystemLinux, instanceOne.OperatingSystem)

	instanceTwo := instances[1]
	assert.Equal(t, "192.0.0.2", instanceTwo.IpAddress)
	assert.Equal(t, "us-east-instance-2", instanceTwo.InstanceId)
	assert.Equal(t, "us-east-1", instanceTwo.Region)
	assert.Equal(t, fleetId, instanceTwo.FleetId)
	assert.Equal(t, config.OperatingSystemWindows, instanceTwo.OperatingSystem)
}

// TestGetInstancesFleetLocationLookupError verifies that we properly handle errors when looking up fleet locations
func TestGetInstancesFleetLocationLookupError(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	expectedErr := errors.New("my error")

	awsMock.DescribeFleetLocationAttributesFunc = func(ctx context.Context, params *gamelift.DescribeFleetLocationAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetLocationAttributesOutput, error) {
		return nil, expectedErr
	}

	_, err := client.GetInstances(context.Background(), fleetId, []string{})

	assert.ErrorIs(t, err, expectedErr)
}

// TestGetInstancesUnsupportedRegionError verifies that we properly handle the case where the call to DescribeFleetLocationAttributes returns an UnsupportedRegionException.
// We should still try to fetch instances for fleets that do not have multi-location support.
func TestGetInstancesUnsupportedRegionError(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	// Return UnsupportedRegionException when DescribeFleetLocationAttributes is called
	awsMock.DescribeFleetLocationAttributesFunc = func(ctx context.Context, params *gamelift.DescribeFleetLocationAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetLocationAttributesOutput, error) {
		err := new(types.UnsupportedRegionException)
		return nil, err
	}

	// DescribeInstances will still function without the location parameter passed
	awsMock.DescribeInstancesFunc = func(ctx context.Context, params *gamelift.DescribeInstancesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeInstancesOutput, error) {
		return &gamelift.DescribeInstancesOutput{
			Instances: []types.Instance{
				{OperatingSystem: types.OperatingSystemAmazonLinux2023, IpAddress: aws.String("127.0.0.2"), InstanceId: aws.String("us-west-instance-1"), Location: aws.String("us-west-1"), Status: types.InstanceStatusActive},
			},
		}, nil
	}

	instances, err := client.GetInstances(context.Background(), fleetId, []string{})

	assert.Nil(t, err)

	// Ensure we don't pass a location to the describe instances call
	assert.Len(t, awsMock.DescribeInstancesCalls(), 1)
	assert.Nil(t, awsMock.DescribeInstancesCalls()[0].Params.Location)

	// Ensure the active instances are still returned
	assert.Len(t, instances, 1)
}

// TestGetInstancesLookupInstancesError verifies that we properly handle errors when looking up instances
func TestGetInstancesLookupInstancesError(t *testing.T) {
	awsMock := &AWSGameliftClientMock{}
	client := &GameLiftClient{gamelift: awsMock}

	expectedErr := errors.New("my error")

	// Set up a fleet with two active regions (will be used), and one activating region (will be ignored)
	awsMock.DescribeFleetLocationAttributesFunc = func(ctx context.Context, params *gamelift.DescribeFleetLocationAttributesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeFleetLocationAttributesOutput, error) {
		return &gamelift.DescribeFleetLocationAttributesOutput{
			LocationAttributes: []types.LocationAttributes{
				types.LocationAttributes{LocationState: &types.LocationState{Location: aws.String("us-west-1"), Status: types.FleetStatusActive}},
			},
		}, nil
	}

	awsMock.DescribeInstancesFunc = func(ctx context.Context, params *gamelift.DescribeInstancesInput, optFns ...func(*gamelift.Options)) (*gamelift.DescribeInstancesOutput, error) {
		return nil, expectedErr
	}

	_, err := client.GetInstances(context.Background(), fleetId, []string{})

	assert.ErrorIs(t, err, expectedErr)
}
