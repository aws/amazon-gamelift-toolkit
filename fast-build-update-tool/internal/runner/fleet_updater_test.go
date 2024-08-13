package runner

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/ssh"
)

var (
	fleetId = "fleet-1234"
)

func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
}

type FleetUpdaterTestSuite struct {
	suite.Suite

	currentDir      string
	buildZipPath    string
	privateKeyPath  string
	defaultArgs     config.CLIArgs
	defaultInstance *gamelift.Instance
}

func TestFleetUpdater(t *testing.T) {
	suite.Run(t, new(FleetUpdaterTestSuite))
}

func (s *FleetUpdaterTestSuite) SetupSuite() {
	_, s.privateKeyPath = generatePrivateSSHKey()
}

func (s *FleetUpdaterTestSuite) SetupTest() {
	s.currentDir, _ = os.Getwd()
	s.buildZipPath = filepath.Join(s.currentDir, "testdata", "game-executable.zip")

	s.defaultArgs = config.CLIArgs{
		FleetId:        fleetId,
		IpRange:        "0.0.0.0/0",
		BuildZipPath:   s.buildZipPath,
		SSHPort:        22,
		InstanceIds:    make([]string, 0),
		RestartProcess: false,
		LockName:       "test",
		Verbose:        false,
		PrivateKeyPath: s.privateKeyPath,
	}

	s.defaultInstance = &gamelift.Instance{
		IpAddress:       "127.0.0.1",
		InstanceId:      "i-12345",
		Region:          "us-east-1",
		OperatingSystem: config.OperatingSystemLinux,
		FleetId:         fleetId,
	}
}

func (s *FleetUpdaterTestSuite) TearDownSuite() {
	if s.privateKeyPath != "" {
		os.Remove(s.privateKeyPath)
	}
}

func generatePrivateSSHKey() (ssh.Signer, string) {
	// Generate a temporary private key file
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	})

	keyFile, _ := os.CreateTemp("", "key")
	defer keyFile.Close()

	_, _ = keyFile.Write(pemBytes)

	signer, _ := ssh.NewSignerFromKey(privateKey)

	return signer, keyFile.Name()
}

// TestUpdateInstancesSuccess run through a normal fleet update operation where everything succeeds
func (s *FleetUpdaterTestSuite) TestUpdateInstancesSuccess() {
	t := s.T()

	logger := NewTestLogger()

	gameliftClient := &GameLiftClientMock{
		GetFleetFunc: func(ctx context.Context, fleetId string) (*gamelift.Fleet, error) {
			return &gamelift.Fleet{Id: fleetId, OperatingSystem: config.OperatingSystemLinux, ExecutablePaths: []string{"bin/server.exe"}}, nil
		},
		GetInstancesFunc: func(ctx context.Context, fleetId string, allowedInstanceIds []string) ([]*gamelift.Instance, error) {
			return []*gamelift.Instance{s.defaultInstance}, nil
		},
		OpenPortForFleetFunc: func(ctx context.Context, fleetId string, port int32, ipRange string) error {
			return nil
		},
	}

	instanceUpdaterFactory := &InstanceUpdaterFactoryMock{
		CreateFunc: func(ctx context.Context, verbose bool, sshKey ssh.Signer, updateScript string, sshPort int32, instance *gamelift.Instance) (InstanceUpdater, error) {
			return &InstanceUpdaterMock{
				UpdateFunc: func(ctx context.Context) error {
					return nil
				},
			}, nil
		},
	}

	f := &FleetUpdater{
		args:                   s.defaultArgs,
		gameLiftClient:         gameliftClient,
		logger:                 logger,
		updateScriptGenerator:  tools.NewInstanceUpdateScriptGenerator(s.defaultArgs.GetUpdateOperation(), s.defaultArgs.BuildZipPath, s.defaultArgs.LockName),
		sshConfigManager:       tools.NewSSHConfigManager(logger, s.defaultArgs.PrivateKeyPath, s.defaultArgs.SSHPort),
		zipValidator:           tools.NewZipValidator(s.defaultArgs.BuildZipPath),
		instanceUpdaterFactory: instanceUpdaterFactory,
		reportWriter:           NewFleetUpdateReportWriter(s.defaultArgs.FleetId, s.defaultArgs.Verbose),
	}
	defer f.Cleanup()

	results, err := f.UpdateInstances(context.Background())

	assert.Nil(t, err)

	assert.Empty(t, results.InstancesFailedUpdate)
	assert.Equal(t, 1, results.InstancesFound)
	assert.Equal(t, 1, results.InstancesUpdated)

	createCalls := instanceUpdaterFactory.CreateCalls()
	assert.Len(t, createCalls, 1)
	assert.False(t, createCalls[0].Verbose)
	assert.NotNil(t, createCalls[0].SshKey)
	assert.NotEmpty(t, createCalls[0].UpdateScript)
	assert.Equal(t, int32(22), createCalls[0].SshPort)
	assert.Equal(t, s.defaultInstance, createCalls[0].Instance)
}

// TestUpdateInstancesFailed ensures we return proper errors, and results when updating an instance in the fleet fails
func (s *FleetUpdaterTestSuite) TestUpdateInstancesFailed() {
	t := s.T()

	logger := NewTestLogger()

	gameliftClient := &GameLiftClientMock{
		GetFleetFunc: func(ctx context.Context, fleetId string) (*gamelift.Fleet, error) {
			return &gamelift.Fleet{Id: fleetId, OperatingSystem: config.OperatingSystemLinux, ExecutablePaths: []string{"bin/server.exe"}}, nil
		},
		GetInstancesFunc: func(ctx context.Context, fleetId string, allowedInstanceIds []string) ([]*gamelift.Instance, error) {
			return []*gamelift.Instance{s.defaultInstance}, nil
		},
		OpenPortForFleetFunc: func(ctx context.Context, fleetId string, port int32, ipRange string) error {
			return nil
		},
	}

	// Set up an instance updater that fails
	instanceUpdaterFactory := &InstanceUpdaterFactoryMock{
		CreateFunc: func(ctx context.Context, verbose bool, sshKey ssh.Signer, updateScript string, sshPort int32, instance *gamelift.Instance) (InstanceUpdater, error) {
			return &InstanceUpdaterMock{
				UpdateFunc: func(ctx context.Context) error {
					return errors.New("failed to update instance")
				},
			}, nil
		},
	}

	f := &FleetUpdater{
		args:                   s.defaultArgs,
		gameLiftClient:         gameliftClient,
		logger:                 logger,
		updateScriptGenerator:  tools.NewInstanceUpdateScriptGenerator(s.defaultArgs.GetUpdateOperation(), s.defaultArgs.BuildZipPath, s.defaultArgs.LockName),
		sshConfigManager:       tools.NewSSHConfigManager(logger, s.defaultArgs.PrivateKeyPath, s.defaultArgs.SSHPort),
		zipValidator:           tools.NewZipValidator(s.defaultArgs.BuildZipPath),
		instanceUpdaterFactory: instanceUpdaterFactory,
		reportWriter:           NewFleetUpdateReportWriter(s.defaultArgs.FleetId, s.defaultArgs.Verbose),
	}
	defer f.Cleanup()

	results, err := f.UpdateInstances(context.Background())

	assert.NotNil(t, err)

	assert.Len(t, results.InstancesFailedUpdate, 1)
	assert.Equal(t, 1, results.InstancesFound)
	assert.Equal(t, 0, results.InstancesUpdated)
}
