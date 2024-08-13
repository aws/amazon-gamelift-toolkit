// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package tools

import (
	"io"
	"sync"
)

// PTYMock is a mock implementation of PTY.
//
//	func TestSomethingThatUsesPTY(t *testing.T) {
//
//		// make and configure a mocked PTY
//		mockedPTY := &PTYMock{
//			CleanupFunc: func()  {
//				panic("mock out the Cleanup method")
//			},
//			ReaderFunc: func() io.Reader {
//				panic("mock out the Reader method")
//			},
//			RunCommandFunc: func(cmd string) error {
//				panic("mock out the RunCommand method")
//			},
//			StartFunc: func(cmdName string, args []string, env []string) error {
//				panic("mock out the Start method")
//			},
//			WaitFunc: func() error {
//				panic("mock out the Wait method")
//			},
//		}
//
//		// use mockedPTY in code that requires PTY
//		// and then make assertions.
//
//	}
type PTYMock struct {
	// CleanupFunc mocks the Cleanup method.
	CleanupFunc func()

	// ReaderFunc mocks the Reader method.
	ReaderFunc func() io.Reader

	// RunCommandFunc mocks the RunCommand method.
	RunCommandFunc func(cmd string) error

	// StartFunc mocks the Start method.
	StartFunc func(cmdName string, args []string, env []string) error

	// WaitFunc mocks the Wait method.
	WaitFunc func() error

	// calls tracks calls to the methods.
	calls struct {
		// Cleanup holds details about calls to the Cleanup method.
		Cleanup []struct {
		}
		// Reader holds details about calls to the Reader method.
		Reader []struct {
		}
		// RunCommand holds details about calls to the RunCommand method.
		RunCommand []struct {
			// Cmd is the cmd argument value.
			Cmd string
		}
		// Start holds details about calls to the Start method.
		Start []struct {
			// CmdName is the cmdName argument value.
			CmdName string
			// Args is the args argument value.
			Args []string
			// Env is the env argument value.
			Env []string
		}
		// Wait holds details about calls to the Wait method.
		Wait []struct {
		}
	}
	lockCleanup    sync.RWMutex
	lockReader     sync.RWMutex
	lockRunCommand sync.RWMutex
	lockStart      sync.RWMutex
	lockWait       sync.RWMutex
}

// Cleanup calls CleanupFunc.
func (mock *PTYMock) Cleanup() {
	if mock.CleanupFunc == nil {
		panic("PTYMock.CleanupFunc: method is nil but PTY.Cleanup was just called")
	}
	callInfo := struct {
	}{}
	mock.lockCleanup.Lock()
	mock.calls.Cleanup = append(mock.calls.Cleanup, callInfo)
	mock.lockCleanup.Unlock()
	mock.CleanupFunc()
}

// CleanupCalls gets all the calls that were made to Cleanup.
// Check the length with:
//
//	len(mockedPTY.CleanupCalls())
func (mock *PTYMock) CleanupCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockCleanup.RLock()
	calls = mock.calls.Cleanup
	mock.lockCleanup.RUnlock()
	return calls
}

// Reader calls ReaderFunc.
func (mock *PTYMock) Reader() io.Reader {
	if mock.ReaderFunc == nil {
		panic("PTYMock.ReaderFunc: method is nil but PTY.Reader was just called")
	}
	callInfo := struct {
	}{}
	mock.lockReader.Lock()
	mock.calls.Reader = append(mock.calls.Reader, callInfo)
	mock.lockReader.Unlock()
	return mock.ReaderFunc()
}

// ReaderCalls gets all the calls that were made to Reader.
// Check the length with:
//
//	len(mockedPTY.ReaderCalls())
func (mock *PTYMock) ReaderCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockReader.RLock()
	calls = mock.calls.Reader
	mock.lockReader.RUnlock()
	return calls
}

// RunCommand calls RunCommandFunc.
func (mock *PTYMock) RunCommand(cmd string) error {
	if mock.RunCommandFunc == nil {
		panic("PTYMock.RunCommandFunc: method is nil but PTY.RunCommand was just called")
	}
	callInfo := struct {
		Cmd string
	}{
		Cmd: cmd,
	}
	mock.lockRunCommand.Lock()
	mock.calls.RunCommand = append(mock.calls.RunCommand, callInfo)
	mock.lockRunCommand.Unlock()
	return mock.RunCommandFunc(cmd)
}

// RunCommandCalls gets all the calls that were made to RunCommand.
// Check the length with:
//
//	len(mockedPTY.RunCommandCalls())
func (mock *PTYMock) RunCommandCalls() []struct {
	Cmd string
} {
	var calls []struct {
		Cmd string
	}
	mock.lockRunCommand.RLock()
	calls = mock.calls.RunCommand
	mock.lockRunCommand.RUnlock()
	return calls
}

// Start calls StartFunc.
func (mock *PTYMock) Start(cmdName string, args []string, env []string) error {
	if mock.StartFunc == nil {
		panic("PTYMock.StartFunc: method is nil but PTY.Start was just called")
	}
	callInfo := struct {
		CmdName string
		Args    []string
		Env     []string
	}{
		CmdName: cmdName,
		Args:    args,
		Env:     env,
	}
	mock.lockStart.Lock()
	mock.calls.Start = append(mock.calls.Start, callInfo)
	mock.lockStart.Unlock()
	return mock.StartFunc(cmdName, args, env)
}

// StartCalls gets all the calls that were made to Start.
// Check the length with:
//
//	len(mockedPTY.StartCalls())
func (mock *PTYMock) StartCalls() []struct {
	CmdName string
	Args    []string
	Env     []string
} {
	var calls []struct {
		CmdName string
		Args    []string
		Env     []string
	}
	mock.lockStart.RLock()
	calls = mock.calls.Start
	mock.lockStart.RUnlock()
	return calls
}

// Wait calls WaitFunc.
func (mock *PTYMock) Wait() error {
	if mock.WaitFunc == nil {
		panic("PTYMock.WaitFunc: method is nil but PTY.Wait was just called")
	}
	callInfo := struct {
	}{}
	mock.lockWait.Lock()
	mock.calls.Wait = append(mock.calls.Wait, callInfo)
	mock.lockWait.Unlock()
	return mock.WaitFunc()
}

// WaitCalls gets all the calls that were made to Wait.
// Check the length with:
//
//	len(mockedPTY.WaitCalls())
func (mock *PTYMock) WaitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockWait.RLock()
	calls = mock.calls.Wait
	mock.lockWait.RUnlock()
	return calls
}
