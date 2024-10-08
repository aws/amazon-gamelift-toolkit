// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package runner

import (
	"context"
	"golang.org/x/crypto/ssh"
	"sync"
)

// FileUploaderMock is a mock implementation of FileUploader.
//
//	func TestSomethingThatUsesFileUploader(t *testing.T) {
//
//		// make and configure a mocked FileUploader
//		mockedFileUploader := &FileUploaderMock{
//			CopyFilesFunc: func(ctx context.Context, remotePublicKey ssh.PublicKey) error {
//				panic("mock out the CopyFiles method")
//			},
//		}
//
//		// use mockedFileUploader in code that requires FileUploader
//		// and then make assertions.
//
//	}
type FileUploaderMock struct {
	// CopyFilesFunc mocks the CopyFiles method.
	CopyFilesFunc func(ctx context.Context, remotePublicKey ssh.PublicKey) error

	// calls tracks calls to the methods.
	calls struct {
		// CopyFiles holds details about calls to the CopyFiles method.
		CopyFiles []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// RemotePublicKey is the remotePublicKey argument value.
			RemotePublicKey ssh.PublicKey
		}
	}
	lockCopyFiles sync.RWMutex
}

// CopyFiles calls CopyFilesFunc.
func (mock *FileUploaderMock) CopyFiles(ctx context.Context, remotePublicKey ssh.PublicKey) error {
	if mock.CopyFilesFunc == nil {
		panic("FileUploaderMock.CopyFilesFunc: method is nil but FileUploader.CopyFiles was just called")
	}
	callInfo := struct {
		Ctx             context.Context
		RemotePublicKey ssh.PublicKey
	}{
		Ctx:             ctx,
		RemotePublicKey: remotePublicKey,
	}
	mock.lockCopyFiles.Lock()
	mock.calls.CopyFiles = append(mock.calls.CopyFiles, callInfo)
	mock.lockCopyFiles.Unlock()
	return mock.CopyFilesFunc(ctx, remotePublicKey)
}

// CopyFilesCalls gets all the calls that were made to CopyFiles.
// Check the length with:
//
//	len(mockedFileUploader.CopyFilesCalls())
func (mock *FileUploaderMock) CopyFilesCalls() []struct {
	Ctx             context.Context
	RemotePublicKey ssh.PublicKey
} {
	var calls []struct {
		Ctx             context.Context
		RemotePublicKey ssh.PublicKey
	}
	mock.lockCopyFiles.RLock()
	calls = mock.calls.CopyFiles
	mock.lockCopyFiles.RUnlock()
	return calls
}
