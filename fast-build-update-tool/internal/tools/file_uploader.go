package tools

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func runCommand(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	cmd.Stderr = config.NewErrorLogger(cmdName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// FileUploader is used to upload one or more files to a remote instance
type FileUploader struct {
	logger                *slog.Logger
	remoteIpAddress       string
	privateKeyPath        string
	remoteUser            config.RemoteUser
	remoteUploadDirectory config.RemoteUploadDirectory
	filesToUpload         []string
	sshPort               int32
	commandRunner         func(cmdName string, args ...string) error
}

const (
	scpCommand = "scp"
)

// NewFileUploader instantiates a new file uploader for the given GameLift instance
func NewFileUploader(logger *slog.Logger, instance *gamelift.Instance, privateKeyPath string, filesToUpload []string, sshPort int32) (*FileUploader, error) {
	result := &FileUploader{
		logger:                logger.With("context", "FileUploader"),
		remoteIpAddress:       instance.IpAddress,
		privateKeyPath:        privateKeyPath,
		remoteUser:            config.RemoteUserForOperatingSystem(instance.OperatingSystem),
		remoteUploadDirectory: config.RemoteUploadDirectoryForOperatingSystem(instance.OperatingSystem),
		filesToUpload:         filesToUpload,
		sshPort:               sshPort,
		commandRunner:         runCommand,
	}

	return result, result.Validate()
}

// Validate that the FileUploader can copy files to the remote instance
func (f *FileUploader) Validate() error {
	return verifyExe(scpCommand)
}

// CopyFiles opens an connection to the remote instance, and copies files up to it
func (f *FileUploader) CopyFiles(ctx context.Context, remotePublicKey ssh.PublicKey) error {
	tempKnownHostsFile, err := f.generateKnownHostsFile(ctx, remotePublicKey)
	if err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(tempKnownHostsFile); err != nil {
			f.logger.Error("error removing temporary known hosts file", "file", tempKnownHostsFile, "error", err)
		}
	}()

	for _, file := range f.filesToUpload {
		if err := f.copyFile(ctx, file, tempKnownHostsFile); err != nil {
			return fmt.Errorf("error uploading file %s to server %w", file, err)
		}
	}

	return nil
}

// generateKnownHostsFile generates a temporary known hosts file with the provided public key, so we can safely SCP files to server
func (f *FileUploader) generateKnownHostsFile(ctx context.Context, remotePublicKey ssh.PublicKey) (string, error) {
	khFile, err := os.CreateTemp("", "known_hosts")
	if err != nil {
		return "", err
	}

	defer func() {
		if err = khFile.Close(); err != nil {
			f.logger.Warn("error closing known hosts file", "error", err)
		}
	}()

	hostPort := net.JoinHostPort(f.remoteIpAddress, fmt.Sprintf("%d", f.sshPort))
	line := knownhosts.Line([]string{hostPort}, remotePublicKey)
	_, err = khFile.Write([]byte(line))
	if err != nil {
		return "", err
	}

	return khFile.Name(), err
}

// copyFile actually runs the scp command to copy a file ot the server
func (f *FileUploader) copyFile(ctx context.Context, file, knownHostsFile string) error {
	f.logger.Debug("copying file to remote instance", "file", file)

	err := f.commandRunner(scpCommand,
		"-o",
		"UserKnownHostsFile="+knownHostsFile,
		"-P",
		fmt.Sprintf("%d", f.sshPort),
		"-i",
		f.privateKeyPath,
		file,
		fmt.Sprintf("%s@%s:%s", f.remoteUser, f.remoteIpAddress, string(f.remoteUploadDirectory)+filepath.Base(file)))
	if err != nil {
		return fmt.Errorf("error copying file to remote instance: %s %w", f.remoteIpAddress, err)
	}

	return nil
}
