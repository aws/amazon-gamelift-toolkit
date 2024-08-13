package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/gamelift"
	"golang.org/x/crypto/ssh"
)

// SSHEnabler is used to enable, and configure SSH on a remote instance.
// SSH is enabled over AWS SSM. We use SSH along with SSM so that we can upload files to the instance.
// The default SSM configuration for GameLIft does not support the SSH proxy flow.
type SSHEnabler struct {
	logger               *slog.Logger
	instance             *gamelift.Instance
	instanceAccessGetter GameLiftInstanceAccessGetter
	clientPublicKey      string
	isNewCommandOutput   func(output string) bool
	commandsToRun        []string
	pty                  PTY
}

// NewSSHEnabler builds a new SSHEnabler for the target instance
func NewSSHEnabler(logger *slog.Logger, instance *gamelift.Instance, instanceAccessGetter GameLiftInstanceAccessGetter, localPublicKey ssh.PublicKey, sshPort int32) (*SSHEnabler, error) {
	localPublicKeyStr := convertPublicKeyToString(localPublicKey)

	pty, err := newPtyCommandRunner()
	if err != nil {
		return nil, err
	}

	var updateCommands []string
	var isNewCommandOutput func(output string) bool

	switch instance.OperatingSystem {

	case config.OperatingSystemWindows:
		updateCommands = windowsSSHEnableCommands(localPublicKeyStr, sshPort)
		isNewCommandOutput = IsNewCommandOutputWindows

	case config.OperatingSystemLinux:
		updateCommands = linuxSSHEnableCommands(localPublicKeyStr)
		isNewCommandOutput = IsNewCommandOutputLinux

	default:
		return nil, config.UnknownOperatingSystemError(fmt.Sprint(instance.OperatingSystem))
	}

	enabler := &SSHEnabler{
		logger:               logger.With("context", "SSHEnabler"),
		pty:                  pty,
		instance:             instance,
		instanceAccessGetter: instanceAccessGetter,
		clientPublicKey:      localPublicKeyStr,
		isNewCommandOutput:   isNewCommandOutput,
		commandsToRun:        updateCommands,
	}

	return enabler, enabler.Validate()
}

func convertPublicKeyToString(key ssh.PublicKey) string {
	return string(bytes.TrimSuffix(ssh.MarshalAuthorizedKey(key), []byte{'\n'}))
}

func (s *SSHEnabler) Validate() error {
	// Verify we have the AWS CLI in the path
	if err := verifyExe(awsCommand); err != nil {
		return err
	}

	// Verify we have the session manager plugin in the path
	if err := verifyExe(sessionManagerCommand); err != nil {
		return err
	}

	return nil
}

// Enable enable SSH on the remote instance
func (s *SSHEnabler) Enable(ctx context.Context) (ssh.PublicKey, error) {
	defer s.pty.Cleanup()

	// Get remote instance access credentials
	accessCredentials, err := s.instanceAccessGetter.GetInstanceAccess(ctx, s.instance.FleetId, s.instance.InstanceId)
	if err != nil {
		return nil, err
	}

	// Add AWS access credential environment variables
	env := os.Environ()
	env = append(env, envVar("AWS_REGION", s.instance.Region))
	env = append(env, envVar("AWS_ACCESS_KEY_ID", accessCredentials.AccessKeyId))
	env = append(env, envVar("AWS_SECRET_ACCESS_KEY", accessCredentials.SecretAccessKey))
	env = append(env, envVar("AWS_SESSION_TOKEN", accessCredentials.SessionToken))

	err = s.pty.Start("aws", []string{"ssm", "start-session", "--target", s.instance.InstanceId}, env)
	if err != nil {
		return nil, err
	}

	// Channel used to let us know when we can send the next command to the SSM session
	commandReady := make(chan int, 1)

	// Channel used to let us know when the SSM session has written out the contents of the remote server's public key
	sshKeyReady := make(chan string, 1)

	// Set up an io.Writer to handle the remote output of the SSM session
	ioWriter := &ptyWriter{
		logger:             s.logger,
		commandReady:       commandReady,
		commandsToAccept:   len(s.commandsToRun),
		sshKeyReady:        sshKeyReady,
		isNewCommandOutput: s.isNewCommandOutput,
		clientPublicKey:    s.clientPublicKey,
	}

	// Start a goroutine to actually send the commands to the remote session
	go func() {
		for i, command := range s.commandsToRun {
			s.logger.Debug("waiting to run ssh enable command", "commandNumber", i)
			<-commandReady
			s.logger.Debug("running ssh enable command", "commandNumber", i)

			err := s.pty.RunCommand(command)
			if err != nil {
				s.logger.Error("error running remote command", "error", err)
				return
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Start a goroutine to copy the output from the SSM session to our writer
	go func() {
		_, err := io.Copy(ioWriter, s.pty.Reader())
		if err != nil && !errors.Is(err, os.ErrClosed) && errors.Is(err, io.EOF) {
			s.logger.Warn("error copying pty commands from remote instance", "err", err)
		}
	}()

	// Wait for the SSM session to finish
	err = s.pty.Wait()
	if err != nil {
		return nil, err
	}

	// Read the remote public SSH key out of the channel, and parse it
	sshKey := <-sshKeyReady
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey))
	if err != nil {
		return nil, fmt.Errorf("error parsing remote public key %s %w", sshKey, err)
	}

	return pubKey, nil
}

// ptyWriter is used to handle writing, and parsing output from the remote SSM session
type ptyWriter struct {
	logger             *slog.Logger
	commandReady       chan int
	sshKeyReady        chan string
	commandsToAccept   int
	commandsRun        int
	isNewCommandOutput func(output string) bool
	clientPublicKey    string
}

func (w *ptyWriter) Write(p []byte) (int, error) {
	terminalOutputStr := string(p)

	// If we have a new command input (eg. `sh $`), let the channel know we can accept the next command.
	if w.isNewCommandOutput(terminalOutputStr) && w.commandsRun < w.commandsToAccept {
		w.commandsRun = w.commandsRun + 1
		w.commandReady <- w.commandsRun
		// Everything has run, close the channel
		if w.commandsToAccept <= w.commandsRun {
			close(w.commandReady)
		}
	}

	// We'll cat the public key file, make sure we capture the output, we need this to connect to the server later on
	match := FindED25519PublicKey(terminalOutputStr)
	if match != "" && !strings.Contains(w.clientPublicKey, match) {
		w.logger.Debug("found server public key", "key", match)
		w.sshKeyReady <- match
	}

	return len(p), nil
}

var (
	publicKeyRegex = regexp.MustCompile("ssh-ed25519 ([A-Za-z0-9+/=]+)")
)

func FindED25519PublicKey(s string) string {
	return publicKeyRegex.FindString(s)
}

func envVar(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

const (
	awsCommand            = "aws"
	sessionManagerCommand = "session-manager-plugin"
)

//go:generate moq -skip-ensure -out ./moq_gamelift_instance_access_getter_test.go . GameLiftInstanceAccessGetter
type GameLiftInstanceAccessGetter interface {
	GetInstanceAccess(ctx context.Context, fleetId string, instanceId string) (*gamelift.InstanceAccessCredentials, error)
}
