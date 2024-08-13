package gamelift

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	"github.com/aws/aws-sdk-go-v2/service/gamelift/types"
)

// Instance represents a single instance a GameLift fleet
type Instance struct {
	// IpAddress the ip address of the instance
	IpAddress string
	// InstanceId the instance id of the instance
	InstanceId string
	// Region the region the instance is running in (us-east-1, etc...)
	Region string
	// FleetId the id of the fleet this instance belongs to
	FleetId string
	// OperatingSystem the operating system that this instance is running
	OperatingSystem config.OperatingSystem
}

// GetInstances will return all active instances in the provided fleet.
// Optionally this function can filter out any instances not found in allowedInstanceIds.
func (g *GameLiftClient) GetInstances(ctx context.Context, fleetId string, allowedInstanceIds []string) ([]*Instance, error) {
	// First we need to fetch locations, as there may be instances outside of the home region
	locations, err := g.getLocations(ctx, fleetId, make([]string, 0, 1), nil)
	if err != nil {
		return make([]*Instance, 0), err
	}

	// Fetch each instance in each location and add it to results
	result := make([]*Instance, 0, 1)
	for _, location := range locations {
		result, err = g.getInstancesInternal(ctx, fleetId, location, result, nil)
		if err != nil {
			return make([]*Instance, 0), err
		}
	}

	// Return the filtered slice of any instances we found
	return g.filterInstances(result, allowedInstanceIds)
}

func (g *GameLiftClient) getLocations(ctx context.Context, fleetId string, locations []string, nextToken *string) ([]string, error) {
	locationAttributesOutput, err := g.gamelift.DescribeFleetLocationAttributes(ctx, &gamelift.DescribeFleetLocationAttributesInput{
		FleetId:   aws.String(fleetId),
		NextToken: nextToken,
	})
	if err != nil {
		return locations, fmt.Errorf("error checking locations for fleet: %w", err)
	}

	for _, locationAttributes := range locationAttributesOutput.LocationAttributes {
		// Filter out anything that is not active
		if strings.EqualFold(string(locationAttributes.LocationState.Status), string(types.FleetStatusActive)) {
			locations = append(locations, *locationAttributes.LocationState.Location)
		}
	}

	// If the results are paginated, fetch the next page
	if locationAttributesOutput.NextToken != nil {
		return g.getLocations(ctx, fleetId, locations, locationAttributesOutput.NextToken)
	}

	return locations, nil
}

func (g *GameLiftClient) getInstancesInternal(ctx context.Context, fleetId, location string, instances []*Instance, nextToken *string) ([]*Instance, error) {
	instanceOutput, err := g.gamelift.DescribeInstances(ctx, &gamelift.DescribeInstancesInput{
		FleetId:   aws.String(fleetId),
		Location:  aws.String(location),
		NextToken: nextToken,
	})
	if err != nil {
		return instances, fmt.Errorf("error describing instances: %w", err)
	}

	for _, instance := range instanceOutput.Instances {
		// Filter out anything that is not active
		if !strings.EqualFold(string(instance.Status), string(types.InstanceStatusActive)) {
			slog.Debug("instance not active, skipping...", "instanceId", *instance.InstanceId, "status", instance.Status)
			continue
		}

		os, err := operatingSystemLookup(instance.OperatingSystem)
		if err != nil {
			return instances, err
		}

		instances = append(instances, &Instance{
			IpAddress:       *instance.IpAddress,
			InstanceId:      *instance.InstanceId,
			Region:          *instance.Location,
			FleetId:         fleetId,
			OperatingSystem: os,
		})
	}

	// If the results are paginated, fetch the next page
	if instanceOutput.NextToken != nil {
		return g.getInstancesInternal(ctx, fleetId, location, instances, instanceOutput.NextToken)
	}

	return instances, err
}

// filterInstances will filter out any instances not in the allow list
func (g *GameLiftClient) filterInstances(instances []*Instance, allowedInstanceIds []string) ([]*Instance, error) {
	// Nothing to filter, break early
	if len(allowedInstanceIds) == 0 {
		return instances, nil
	}

	// Loop through allow list, and add any matches to the result
	result := make([]*Instance, 0, len(allowedInstanceIds))
	for _, instance := range instances {
		for _, allowedInstanceId := range allowedInstanceIds {
			if instance.InstanceId == allowedInstanceId {
				result = append(result, instance)
			}
		}
	}

	if len(result) != len(allowedInstanceIds) {
		return result, fmt.Errorf("one or more instance ids not found in fleet: %s", strings.Join(allowedInstanceIds, ","))
	}

	return result, nil
}
