// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package tools

import (
	"context"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
	"sync"
)

// GameLiftInstanceGetterMock is a mock implementation of GameLiftInstanceGetter.
//
//	func TestSomethingThatUsesGameLiftInstanceGetter(t *testing.T) {
//
//		// make and configure a mocked GameLiftInstanceGetter
//		mockedGameLiftInstanceGetter := &GameLiftInstanceGetterMock{
//			GetInstancesFunc: func(ctx context.Context, fleetId string) ([]*gamelift.Instance, error) {
//				panic("mock out the GetInstances method")
//			},
//		}
//
//		// use mockedGameLiftInstanceGetter in code that requires GameLiftInstanceGetter
//		// and then make assertions.
//
//	}
type GameLiftInstanceGetterMock struct {
	// GetInstancesFunc mocks the GetInstances method.
	GetInstancesFunc func(ctx context.Context, fleetId string) ([]*gamelift.Instance, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetInstances holds details about calls to the GetInstances method.
		GetInstances []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// FleetId is the fleetId argument value.
			FleetId string
		}
	}
	lockGetInstances sync.RWMutex
}

// GetInstances calls GetInstancesFunc.
func (mock *GameLiftInstanceGetterMock) GetInstances(ctx context.Context, fleetId string) ([]*gamelift.Instance, error) {
	if mock.GetInstancesFunc == nil {
		panic("GameLiftInstanceGetterMock.GetInstancesFunc: method is nil but GameLiftInstanceGetter.GetInstances was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		FleetId string
	}{
		Ctx:     ctx,
		FleetId: fleetId,
	}
	mock.lockGetInstances.Lock()
	mock.calls.GetInstances = append(mock.calls.GetInstances, callInfo)
	mock.lockGetInstances.Unlock()
	return mock.GetInstancesFunc(ctx, fleetId)
}

// GetInstancesCalls gets all the calls that were made to GetInstances.
// Check the length with:
//
//	len(mockedGameLiftInstanceGetter.GetInstancesCalls())
func (mock *GameLiftInstanceGetterMock) GetInstancesCalls() []struct {
	Ctx     context.Context
	FleetId string
} {
	var calls []struct {
		Ctx     context.Context
		FleetId string
	}
	mock.lockGetInstances.RLock()
	calls = mock.calls.GetInstances
	mock.lockGetInstances.RUnlock()
	return calls
}