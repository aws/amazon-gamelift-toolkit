package tools

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
	"github.com/stretchr/testify/assert"
)

// TestCopyFiles verifies we pass the proper arguments to SCP
func TestCopyFiles(t *testing.T) {
	port := 1026
	keyPath := "my.pem"
	uploadFile := "myfile.txt"

	uploader, err := NewFileUploader(NewTestLogger(), &gamelift.Instance{OperatingSystem: config.OperatingSystemLinux, IpAddress: "127.0.0.1"}, keyPath, []string{uploadFile}, int32(port))
	assert.Nil(t, err)

	commandCallCount := 0
	var calledCommand string
	var calledArgs []string

	uploader.commandRunner = func(cmdName string, args ...string) error {
		commandCallCount = commandCallCount + 1
		calledCommand = cmdName
		calledArgs = args
		return nil
	}

	key := testGenerateKey(t)

	err = uploader.CopyFiles(context.Background(), key)
	assert.Nil(t, err)

	assert.Equal(t, 1, commandCallCount)
	assert.Equal(t, "scp", calledCommand)
	assert.Equal(t, "-o", calledArgs[0])
	assert.Contains(t, calledArgs[1], "UserKnownHostsFile=")
	assert.Equal(t, "-P", calledArgs[2])
	assert.Equal(t, strconv.Itoa(port), calledArgs[3])
	assert.Equal(t, "-i", calledArgs[4])
	assert.Equal(t, keyPath, calledArgs[5])
	assert.Equal(t, uploadFile, calledArgs[6])
	assert.Equal(t, "gl-user-remote@127.0.0.1:/tmp/myfile.txt", calledArgs[7])
}

// TestCopyFilesUploadError verifies that we handle any file upload errors properly
func TestCopyFilesUploadError(t *testing.T) {
	uploader, err := NewFileUploader(NewTestLogger(), &gamelift.Instance{OperatingSystem: config.OperatingSystemLinux, IpAddress: "127.0.0.1"}, "mykey.pem", []string{"myfile.txt"}, 1026)
	assert.Nil(t, err)

	expectedErr := errors.New("test error")

	uploader.commandRunner = func(cmdName string, args ...string) error {
		return expectedErr
	}

	key := testGenerateKey(t)

	err = uploader.CopyFiles(context.Background(), key)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, expectedErr.Error())
}
