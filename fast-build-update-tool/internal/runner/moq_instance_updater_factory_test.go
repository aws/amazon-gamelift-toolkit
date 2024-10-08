// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package runner

import (
	"context"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
	"golang.org/x/crypto/ssh"
	"sync"
)

// InstanceUpdaterFactoryMock is a mock implementation of InstanceUpdaterFactory.
//
//	func TestSomethingThatUsesInstanceUpdaterFactory(t *testing.T) {
//
//		// make and configure a mocked InstanceUpdaterFactory
//		mockedInstanceUpdaterFactory := &InstanceUpdaterFactoryMock{
//			CreateFunc: func(ctx context.Context, verbose bool, sshKey ssh.Signer, updateScript string, sshPort int32, instance *gamelift.Instance) (InstanceUpdater, error) {
//				panic("mock out the Create method")
//			},
//		}
//
//		// use mockedInstanceUpdaterFactory in code that requires InstanceUpdaterFactory
//		// and then make assertions.
//
//	}
type InstanceUpdaterFactoryMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, verbose bool, sshKey ssh.Signer, updateScript string, sshPort int32, instance *gamelift.Instance) (InstanceUpdater, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Verbose is the verbose argument value.
			Verbose bool
			// SshKey is the sshKey argument value.
			SshKey ssh.Signer
			// UpdateScript is the updateScript argument value.
			UpdateScript string
			// SshPort is the sshPort argument value.
			SshPort int32
			// Instance is the instance argument value.
			Instance *gamelift.Instance
		}
	}
	lockCreate sync.RWMutex
}

// Create calls CreateFunc.
func (mock *InstanceUpdaterFactoryMock) Create(ctx context.Context, verbose bool, sshKey ssh.Signer, updateScript string, sshPort int32, instance *gamelift.Instance) (InstanceUpdater, error) {
	if mock.CreateFunc == nil {
		panic("InstanceUpdaterFactoryMock.CreateFunc: method is nil but InstanceUpdaterFactory.Create was just called")
	}
	callInfo := struct {
		Ctx          context.Context
		Verbose      bool
		SshKey       ssh.Signer
		UpdateScript string
		SshPort      int32
		Instance     *gamelift.Instance
	}{
		Ctx:          ctx,
		Verbose:      verbose,
		SshKey:       sshKey,
		UpdateScript: updateScript,
		SshPort:      sshPort,
		Instance:     instance,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(ctx, verbose, sshKey, updateScript, sshPort, instance)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedInstanceUpdaterFactory.CreateCalls())
func (mock *InstanceUpdaterFactoryMock) CreateCalls() []struct {
	Ctx          context.Context
	Verbose      bool
	SshKey       ssh.Signer
	UpdateScript string
	SshPort      int32
	Instance     *gamelift.Instance
} {
	var calls []struct {
		Ctx          context.Context
		Verbose      bool
		SshKey       ssh.Signer
		UpdateScript string
		SshPort      int32
		Instance     *gamelift.Instance
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}
