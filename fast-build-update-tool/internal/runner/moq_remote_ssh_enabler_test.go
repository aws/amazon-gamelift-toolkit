// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package runner

import (
	"context"
	"golang.org/x/crypto/ssh"
	"sync"
)

// RemoteSSHEnablerMock is a mock implementation of RemoteSSHEnabler.
//
//	func TestSomethingThatUsesRemoteSSHEnabler(t *testing.T) {
//
//		// make and configure a mocked RemoteSSHEnabler
//		mockedRemoteSSHEnabler := &RemoteSSHEnablerMock{
//			EnableFunc: func(ctx context.Context) (ssh.PublicKey, error) {
//				panic("mock out the Enable method")
//			},
//		}
//
//		// use mockedRemoteSSHEnabler in code that requires RemoteSSHEnabler
//		// and then make assertions.
//
//	}
type RemoteSSHEnablerMock struct {
	// EnableFunc mocks the Enable method.
	EnableFunc func(ctx context.Context) (ssh.PublicKey, error)

	// calls tracks calls to the methods.
	calls struct {
		// Enable holds details about calls to the Enable method.
		Enable []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
	}
	lockEnable sync.RWMutex
}

// Enable calls EnableFunc.
func (mock *RemoteSSHEnablerMock) Enable(ctx context.Context) (ssh.PublicKey, error) {
	if mock.EnableFunc == nil {
		panic("RemoteSSHEnablerMock.EnableFunc: method is nil but RemoteSSHEnabler.Enable was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockEnable.Lock()
	mock.calls.Enable = append(mock.calls.Enable, callInfo)
	mock.lockEnable.Unlock()
	return mock.EnableFunc(ctx)
}

// EnableCalls gets all the calls that were made to Enable.
// Check the length with:
//
//	len(mockedRemoteSSHEnabler.EnableCalls())
func (mock *RemoteSSHEnablerMock) EnableCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockEnable.RLock()
	calls = mock.calls.Enable
	mock.lockEnable.RUnlock()
	return calls
}
