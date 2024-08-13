package tools

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"golang.org/x/crypto/ssh"
)

// SSHConfigManager manages configuration around interactions with SSH
type SSHConfigManager struct {
	logger         *slog.Logger
	privateKeyPath string
	sshPort        int
}

// NewSSHConfigManager builds a new SSHConfigManager
func NewSSHConfigManager(logger *slog.Logger, privateKeyPath string, sshPort int) *SSHConfigManager {
	return &SSHConfigManager{
		logger:         logger.With("context", "LocalSSHConfigManager"),
		privateKeyPath: privateKeyPath,
		sshPort:        sshPort,
	}
}

// DeterminePort will return the proper SSH port for the provided operating system, and validate it against user arguments
func (s *SSHConfigManager) DeterminePort(operatingSystem config.OperatingSystem) (int32, error) {
	port := int32(s.sshPort)

	// Currently this tool does not support using a custom port for Linux
	if operatingSystem == config.OperatingSystemLinux {
		if port > 0 && port != config.DefaultPortLinux {
			s.logger.Warn("custom SSH ports are not supported for Linux fleets, using the default", "port", config.DefaultPortLinux)
		}
		return config.DefaultPortLinux, nil
	}

	// If the user provided a port, but it is invalid return an error
	if port > 0 && port < config.DefaultPortWindows {
		return 0, fmt.Errorf("ssh port must be greater than or equal to %d for Windows servers", config.DefaultPortWindows)
	}

	if port == 0 {
		return config.DefaultPortWindows, nil
	}

	return port, nil
}

// LoadKey loads an SSH key off of the filesystem
func (s *SSHConfigManager) LoadKey(ctx context.Context) (signer ssh.Signer, err error) {
	privateKeyBytes, err := os.ReadFile(s.privateKeyPath)
	if err != nil {
		return signer, fmt.Errorf("error reading private key file for instance: %w", err)
	}

	signer, err = ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return signer, fmt.Errorf("error parsing private key file %w", err)
	}

	return signer, nil
}
