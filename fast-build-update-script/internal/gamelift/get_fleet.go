package gamelift

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	"github.com/aws/aws-sdk-go-v2/service/gamelift/types"
)

// Fleet represents a single GameLift fleet
type Fleet struct {
	// Id of the GameLift fleet
	Id string
	// OperatingSystem that all instances in the fleet run on
	OperatingSystem config.OperatingSystem
	// ExecutablePaths is a slice of all executable paths described by the runtime configuration for this fleet
	ExecutablePaths []string
}

// GetFleet will look up any fleet attributes relevant to this application
func (g *GameLiftClient) GetFleet(ctx context.Context, fleetId string) (*Fleet, error) {
	fleetAttributesOutput, err := g.gamelift.DescribeFleetAttributes(ctx, &gamelift.DescribeFleetAttributesInput{
		FleetIds: []string{fleetId},
	})
	if err != nil {
		return nil, fmt.Errorf("error looking up fleet attributes %w", err)
	}

	if len(fleetAttributesOutput.FleetAttributes) == 0 {
		return nil, fmt.Errorf("fleet not found: %s", fleetId)
	}

	os, err := operatingSystemLookup(fleetAttributesOutput.FleetAttributes[0].OperatingSystem)
	if err != nil {
		return nil, err
	}

	runtimeConfigurationOutput, err := g.gamelift.DescribeRuntimeConfiguration(ctx, &gamelift.DescribeRuntimeConfigurationInput{
		FleetId: aws.String(fleetId),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting fleet runtime configuration %w", err)
	}

	executablePaths := make([]string, 0, len(runtimeConfigurationOutput.RuntimeConfiguration.ServerProcesses))
	for _, proc := range runtimeConfigurationOutput.RuntimeConfiguration.ServerProcesses {
		if !contains(executablePaths, *proc.LaunchPath) {
			executablePaths = append(executablePaths, *proc.LaunchPath)
		}
	}

	return &Fleet{
		Id:              fleetId,
		OperatingSystem: os,
		ExecutablePaths: executablePaths,
	}, nil
}

func contains(values []string, lookup string) bool {
	for _, v := range values {
		if v == lookup {
			return true
		}
	}
	return false
}

func operatingSystemLookup(os types.OperatingSystem) (config.OperatingSystem, error) {
	remoteOSString := strings.ToLower(string(os))

	if strings.Contains(remoteOSString, "windows") {
		return config.OperatingSystemWindows, nil

	} else if strings.Contains(remoteOSString, "linux") {
		return config.OperatingSystemLinux, nil

	} else {
		return config.OperatingSystemUnknown, config.UnknownOperatingSystemError(remoteOSString)
	}
}
