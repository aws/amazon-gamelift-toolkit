// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package runner

import (
	"context"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
	"sync"
)

// GameLiftClientMock is a mock implementation of GameLiftClient.
//
//	func TestSomethingThatUsesGameLiftClient(t *testing.T) {
//
//		// make and configure a mocked GameLiftClient
//		mockedGameLiftClient := &GameLiftClientMock{
//			GetFleetFunc: func(ctx context.Context, fleetId string) (*gamelift.Fleet, error) {
//				panic("mock out the GetFleet method")
//			},
//			GetInstanceAccessFunc: func(ctx context.Context, fleetId string, instanceId string) (*gamelift.InstanceAccessCredentials, error) {
//				panic("mock out the GetInstanceAccess method")
//			},
//			GetInstancesFunc: func(ctx context.Context, fleetId string, allowedInstanceIds []string) ([]*gamelift.Instance, error) {
//				panic("mock out the GetInstances method")
//			},
//			OpenPortForFleetFunc: func(ctx context.Context, fleetId string, port int32, ipRange string) error {
//				panic("mock out the OpenPortForFleet method")
//			},
//		}
//
//		// use mockedGameLiftClient in code that requires GameLiftClient
//		// and then make assertions.
//
//	}
type GameLiftClientMock struct {
	// GetFleetFunc mocks the GetFleet method.
	GetFleetFunc func(ctx context.Context, fleetId string) (*gamelift.Fleet, error)

	// GetInstanceAccessFunc mocks the GetInstanceAccess method.
	GetInstanceAccessFunc func(ctx context.Context, fleetId string, instanceId string) (*gamelift.InstanceAccessCredentials, error)

	// GetInstancesFunc mocks the GetInstances method.
	GetInstancesFunc func(ctx context.Context, fleetId string, allowedInstanceIds []string) ([]*gamelift.Instance, error)

	// OpenPortForFleetFunc mocks the OpenPortForFleet method.
	OpenPortForFleetFunc func(ctx context.Context, fleetId string, port int32, ipRange string) error

	// calls tracks calls to the methods.
	calls struct {
		// GetFleet holds details about calls to the GetFleet method.
		GetFleet []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// FleetId is the fleetId argument value.
			FleetId string
		}
		// GetInstanceAccess holds details about calls to the GetInstanceAccess method.
		GetInstanceAccess []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// FleetId is the fleetId argument value.
			FleetId string
			// InstanceId is the instanceId argument value.
			InstanceId string
		}
		// GetInstances holds details about calls to the GetInstances method.
		GetInstances []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// FleetId is the fleetId argument value.
			FleetId string
			// AllowedInstanceIds is the allowedInstanceIds argument value.
			AllowedInstanceIds []string
		}
		// OpenPortForFleet holds details about calls to the OpenPortForFleet method.
		OpenPortForFleet []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// FleetId is the fleetId argument value.
			FleetId string
			// Port is the port argument value.
			Port int32
			// IpRange is the ipRange argument value.
			IpRange string
		}
	}
	lockGetFleet          sync.RWMutex
	lockGetInstanceAccess sync.RWMutex
	lockGetInstances      sync.RWMutex
	lockOpenPortForFleet  sync.RWMutex
}

// GetFleet calls GetFleetFunc.
func (mock *GameLiftClientMock) GetFleet(ctx context.Context, fleetId string) (*gamelift.Fleet, error) {
	if mock.GetFleetFunc == nil {
		panic("GameLiftClientMock.GetFleetFunc: method is nil but GameLiftClient.GetFleet was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		FleetId string
	}{
		Ctx:     ctx,
		FleetId: fleetId,
	}
	mock.lockGetFleet.Lock()
	mock.calls.GetFleet = append(mock.calls.GetFleet, callInfo)
	mock.lockGetFleet.Unlock()
	return mock.GetFleetFunc(ctx, fleetId)
}

// GetFleetCalls gets all the calls that were made to GetFleet.
// Check the length with:
//
//	len(mockedGameLiftClient.GetFleetCalls())
func (mock *GameLiftClientMock) GetFleetCalls() []struct {
	Ctx     context.Context
	FleetId string
} {
	var calls []struct {
		Ctx     context.Context
		FleetId string
	}
	mock.lockGetFleet.RLock()
	calls = mock.calls.GetFleet
	mock.lockGetFleet.RUnlock()
	return calls
}

// GetInstanceAccess calls GetInstanceAccessFunc.
func (mock *GameLiftClientMock) GetInstanceAccess(ctx context.Context, fleetId string, instanceId string) (*gamelift.InstanceAccessCredentials, error) {
	if mock.GetInstanceAccessFunc == nil {
		panic("GameLiftClientMock.GetInstanceAccessFunc: method is nil but GameLiftClient.GetInstanceAccess was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		FleetId    string
		InstanceId string
	}{
		Ctx:        ctx,
		FleetId:    fleetId,
		InstanceId: instanceId,
	}
	mock.lockGetInstanceAccess.Lock()
	mock.calls.GetInstanceAccess = append(mock.calls.GetInstanceAccess, callInfo)
	mock.lockGetInstanceAccess.Unlock()
	return mock.GetInstanceAccessFunc(ctx, fleetId, instanceId)
}

// GetInstanceAccessCalls gets all the calls that were made to GetInstanceAccess.
// Check the length with:
//
//	len(mockedGameLiftClient.GetInstanceAccessCalls())
func (mock *GameLiftClientMock) GetInstanceAccessCalls() []struct {
	Ctx        context.Context
	FleetId    string
	InstanceId string
} {
	var calls []struct {
		Ctx        context.Context
		FleetId    string
		InstanceId string
	}
	mock.lockGetInstanceAccess.RLock()
	calls = mock.calls.GetInstanceAccess
	mock.lockGetInstanceAccess.RUnlock()
	return calls
}

// GetInstances calls GetInstancesFunc.
func (mock *GameLiftClientMock) GetInstances(ctx context.Context, fleetId string, allowedInstanceIds []string) ([]*gamelift.Instance, error) {
	if mock.GetInstancesFunc == nil {
		panic("GameLiftClientMock.GetInstancesFunc: method is nil but GameLiftClient.GetInstances was just called")
	}
	callInfo := struct {
		Ctx                context.Context
		FleetId            string
		AllowedInstanceIds []string
	}{
		Ctx:                ctx,
		FleetId:            fleetId,
		AllowedInstanceIds: allowedInstanceIds,
	}
	mock.lockGetInstances.Lock()
	mock.calls.GetInstances = append(mock.calls.GetInstances, callInfo)
	mock.lockGetInstances.Unlock()
	return mock.GetInstancesFunc(ctx, fleetId, allowedInstanceIds)
}

// GetInstancesCalls gets all the calls that were made to GetInstances.
// Check the length with:
//
//	len(mockedGameLiftClient.GetInstancesCalls())
func (mock *GameLiftClientMock) GetInstancesCalls() []struct {
	Ctx                context.Context
	FleetId            string
	AllowedInstanceIds []string
} {
	var calls []struct {
		Ctx                context.Context
		FleetId            string
		AllowedInstanceIds []string
	}
	mock.lockGetInstances.RLock()
	calls = mock.calls.GetInstances
	mock.lockGetInstances.RUnlock()
	return calls
}

// OpenPortForFleet calls OpenPortForFleetFunc.
func (mock *GameLiftClientMock) OpenPortForFleet(ctx context.Context, fleetId string, port int32, ipRange string) error {
	if mock.OpenPortForFleetFunc == nil {
		panic("GameLiftClientMock.OpenPortForFleetFunc: method is nil but GameLiftClient.OpenPortForFleet was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		FleetId string
		Port    int32
		IpRange string
	}{
		Ctx:     ctx,
		FleetId: fleetId,
		Port:    port,
		IpRange: ipRange,
	}
	mock.lockOpenPortForFleet.Lock()
	mock.calls.OpenPortForFleet = append(mock.calls.OpenPortForFleet, callInfo)
	mock.lockOpenPortForFleet.Unlock()
	return mock.OpenPortForFleetFunc(ctx, fleetId, port, ipRange)
}

// OpenPortForFleetCalls gets all the calls that were made to OpenPortForFleet.
// Check the length with:
//
//	len(mockedGameLiftClient.OpenPortForFleetCalls())
func (mock *GameLiftClientMock) OpenPortForFleetCalls() []struct {
	Ctx     context.Context
	FleetId string
	Port    int32
	IpRange string
} {
	var calls []struct {
		Ctx     context.Context
		FleetId string
		Port    int32
		IpRange string
	}
	mock.lockOpenPortForFleet.RLock()
	calls = mock.calls.OpenPortForFleet
	mock.lockOpenPortForFleet.RUnlock()
	return calls
}
