package runner

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/gamelift"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/ssh"
)

const (
	instanceId = "i-1234"
)

type InstanceUpdaterTestSuite struct {
	suite.Suite

	publicKey       ssh.PublicKey
	progressTracker *InstanceProgressWriter

	sshEnabler    *RemoteSSHEnablerMock
	fileUploader  *FileUploaderMock
	commandRunner *CommandRunnerMock
}

func TestInstanceUpdater(t *testing.T) {
	suite.Run(t, new(InstanceUpdaterTestSuite))
}

func (s *InstanceUpdaterTestSuite) SetupTest() {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	s.progressTracker, _ = NewInstanceProgressWriter(&gamelift.Instance{InstanceId: instanceId, IpAddress: "127.0.0.1"}, false)

	s.publicKey, _ = ssh.NewPublicKey(&privateKey.PublicKey)

	s.sshEnabler = &RemoteSSHEnablerMock{
		EnableFunc: func(ctx context.Context) (ssh.PublicKey, error) {
			return s.publicKey, nil
		},
	}

	s.fileUploader = &FileUploaderMock{
		CopyFilesFunc: func(ctx context.Context, remotePublicKey ssh.PublicKey) error {
			return nil
		},
	}

	s.commandRunner = &CommandRunnerMock{
		RunFunc: func(ctx context.Context, remotePublicKey ssh.PublicKey) error {
			return nil
		},
	}
}

// TestInstanceUpdate verifies that the golden path of Update() functions correctly and returns valid results
func (s *InstanceUpdaterTestSuite) TestInstanceUpdate() {
	t := s.T()

	updater := &instanceUpdater{
		progressTracker: s.progressTracker,
		sshEnabler:      s.sshEnabler,
		fileUploader:    s.fileUploader,
		commandRunner:   s.commandRunner,
		logger:          NewTestLogger(),
	}

	err := updater.Update(context.Background())
	assert.Nil(t, err)

	assert.Len(t, s.sshEnabler.EnableCalls(), 1)

	assert.Len(t, s.fileUploader.CopyFilesCalls(), 1)
	assert.Equal(t, s.publicKey, s.fileUploader.CopyFilesCalls()[0].RemotePublicKey)

	assert.Len(t, s.commandRunner.RunCalls(), 1)
	assert.Equal(t, s.publicKey, s.commandRunner.RunCalls()[0].RemotePublicKey)
}

// TestInstanceUpdate verifies that enabling ssh shortcuts the process and returns the proper error
func (s *InstanceUpdaterTestSuite) TestInstanceEnableSSHFail() {
	t := s.T()

	expectedErr := errors.New("enable fail")

	s.sshEnabler = &RemoteSSHEnablerMock{
		EnableFunc: func(ctx context.Context) (ssh.PublicKey, error) {
			return nil, expectedErr
		},
	}

	updater := &instanceUpdater{
		progressTracker: s.progressTracker,
		sshEnabler:      s.sshEnabler,
		fileUploader:    s.fileUploader,
		commandRunner:   s.commandRunner,
		logger:          NewTestLogger(),
	}

	err := updater.Update(context.Background())
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, expectedErr.Error())

	assert.Len(t, s.sshEnabler.EnableCalls(), 1)
	assert.Len(t, s.fileUploader.CopyFilesCalls(), 0)
	assert.Len(t, s.commandRunner.RunCalls(), 0)
}

// TestInstanceUpdate verifies that failing to uploading file to the remote instance shortcuts the update process and returns the proper error
func (s *InstanceUpdaterTestSuite) TestInstanceUploadFileFail() {
	t := s.T()

	expectedErr := errors.New("enable fail")

	s.fileUploader = &FileUploaderMock{
		CopyFilesFunc: func(ctx context.Context, remotePublicKey ssh.PublicKey) error {
			return expectedErr
		},
	}

	updater := &instanceUpdater{
		progressTracker: s.progressTracker,
		sshEnabler:      s.sshEnabler,
		fileUploader:    s.fileUploader,
		commandRunner:   s.commandRunner,
		logger:          NewTestLogger(),
	}

	err := updater.Update(context.Background())
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, expectedErr.Error())

	assert.Len(t, s.sshEnabler.EnableCalls(), 1)

	assert.Len(t, s.fileUploader.CopyFilesCalls(), 1)
	assert.Equal(t, s.publicKey, s.fileUploader.CopyFilesCalls()[0].RemotePublicKey)

	assert.Len(t, s.commandRunner.RunCalls(), 0)
}

// TestInstanceRunCommandFail verifies that when running remote commands on the instance fails, the proper error is returned
func (s *InstanceUpdaterTestSuite) TestInstanceRunCommandFail() {
	t := s.T()

	expectedErr := errors.New("enable fail")

	s.commandRunner = &CommandRunnerMock{
		RunFunc: func(ctx context.Context, remotePublicKey ssh.PublicKey) error {
			return expectedErr
		},
	}

	updater := &instanceUpdater{
		progressTracker: s.progressTracker,
		sshEnabler:      s.sshEnabler,
		fileUploader:    s.fileUploader,
		commandRunner:   s.commandRunner,
		logger:          NewTestLogger(),
	}

	err := updater.Update(context.Background())
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, expectedErr.Error())

	assert.Len(t, s.sshEnabler.EnableCalls(), 1)

	assert.Len(t, s.fileUploader.CopyFilesCalls(), 1)
	assert.Equal(t, s.publicKey, s.fileUploader.CopyFilesCalls()[0].RemotePublicKey)

	assert.Len(t, s.commandRunner.RunCalls(), 1)
	assert.Equal(t, s.publicKey, s.commandRunner.RunCalls()[0].RemotePublicKey)
}
