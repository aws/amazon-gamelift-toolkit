package runner

import (
	"context"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
)

//go:generate moq -skip-ensure -out ./moq_gamelift_client_test.go . GameLiftClient

type GameLiftClient interface {
	GetFleet(ctx context.Context, fleetId string) (*gamelift.Fleet, error)
	GetInstanceAccess(ctx context.Context, fleetId string, instanceId string) (*gamelift.InstanceAccessCredentials, error)
	GetInstances(ctx context.Context, fleetId string, allowedInstanceIds []string) ([]*gamelift.Instance, error)
	OpenPortForFleet(ctx context.Context, fleetId string, port int32, ipRange string) error
}

type FleetUpdateResults struct {
	InstancesFound        int
	InstancesUpdated      int
	InstancesFailedUpdate []string
}

type InstanceUpdateState uint

const (
	UpdateStateNotStarted      InstanceUpdateState = iota
	UpdateStateEnableSSH       InstanceUpdateState = iota
	UpdateStateCopyBuild       InstanceUpdateState = iota
	UpdateStateRunUpdateScript InstanceUpdateState = iota

	// Must be last

	UpdateStateCount InstanceUpdateState = iota
)

func (i InstanceUpdateState) String() string {
	switch i {
	case UpdateStateEnableSSH:
		return "enabling remote access"
	case UpdateStateCopyBuild:
		return "copying build to instance"
	case UpdateStateRunUpdateScript:
		return "updating instance"
	case UpdateStateCount:
		return "done"
	default:
		return "not started"
	}
}
