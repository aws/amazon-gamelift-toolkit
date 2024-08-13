package runner

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/ssh"
)

//go:generate moq -skip-ensure -out ./moq_remote_ssh_enabler_test.go . RemoteSSHEnabler
//go:generate moq -skip-ensure -out ./moq_remote_command_runner_test.go . CommandRunner
//go:generate moq -skip-ensure -out ./moq_file_uploader_test.go . FileUploader
//go:generate moq -skip-ensure -out ./moq_instance_updater_test.go . InstanceUpdater

// RemoteSSHEnabler is an abstraction around enabling access to an instance over SSH
type RemoteSSHEnabler interface {
	// Enable SSH on the remote instance
	Enable(ctx context.Context) (remotePublicKey ssh.PublicKey, err error)
}

// CommandRunner is an abstraction around running commands on a remote instance
type CommandRunner interface {
	// Run the command provided on the remote instance
	Run(ctx context.Context, remotePublicKey ssh.PublicKey) error
}

// FileUploader is an abstraction around copying files to a remote instance
type FileUploader interface {
	// CopyFiles will copy files to the remote instance
	CopyFiles(ctx context.Context, remotePublicKey ssh.PublicKey) error
}

// InstanceUpdater is used to update a single instance in a GameLift fleet
type InstanceUpdater interface {
	// Update will trigger the update process for a single instance
	Update(ctx context.Context) error
}

type instanceUpdater struct {
	progressTracker *InstanceProgressWriter
	sshEnabler      RemoteSSHEnabler
	fileUploader    FileUploader
	commandRunner   CommandRunner

	logger *slog.Logger
}

func (s *instanceUpdater) Update(ctx context.Context) error {
	remotePublicKey, err := s.enableSSH(ctx)
	if err != nil {
		return s.processError(err)
	}

	err = s.copyFilesToRemoteInstance(ctx, remotePublicKey)
	if err != nil {
		return s.processError(err)
	}

	err = s.runUpdateScript(ctx, remotePublicKey)
	if err != nil {
		return s.processError(err)
	}

	s.progressTracker.UpdateState(UpdateStateCount)

	return nil
}

func (s *instanceUpdater) processError(err error) error {
	s.progressTracker.UpdateFailed(err)
	return err
}

// enableSSH will enable SSH on the instance. This must happen first as the other Update steps all depend on it.
func (s *instanceUpdater) enableSSH(ctx context.Context) (ssh.PublicKey, error) {
	s.logger.Debug("enabling ssh on remote instance")

	s.progressTracker.UpdateState(UpdateStateEnableSSH)

	remotePublicKey, err := s.sshEnabler.Enable(ctx)
	if err != nil {
		return nil, fmt.Errorf("error enabling ssh on remote instance %w", err)
	}

	s.logger.Debug("done enabling ssh on remote instance")

	return remotePublicKey, nil
}

// copyFilesToRemoteInstance will copy the build and any relevant update scripts to the instance
func (s *instanceUpdater) copyFilesToRemoteInstance(ctx context.Context, remotePublicKey ssh.PublicKey) error {
	s.logger.Debug("copying files to remote instance")

	s.progressTracker.UpdateState(UpdateStateCopyBuild)

	err := s.fileUploader.CopyFiles(ctx, remotePublicKey)
	if err != nil {
		return fmt.Errorf("error copying files to remote instance %w", err)
	}

	s.logger.Debug("done copying files to remote instance")

	return nil
}

// runUpdateScript will actually run a script on the instance to perform the update
func (s *instanceUpdater) runUpdateScript(ctx context.Context, remotePublicKey ssh.PublicKey) error {
	s.logger.Debug("running update script")

	s.progressTracker.UpdateState(UpdateStateRunUpdateScript)

	err := s.commandRunner.Run(ctx, remotePublicKey)
	if err != nil {
		return fmt.Errorf("error running remote command %w", err)
	}

	s.logger.Debug("done running update script")

	return nil
}
