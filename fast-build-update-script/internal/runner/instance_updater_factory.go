package runner

import (
	"context"
	"log/slog"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/gamelift"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/tools"
	"golang.org/x/crypto/ssh"
)

//go:generate moq -skip-ensure -out ./moq_instance_updater_factory_test.go . InstanceUpdaterFactory

// InstanceUpdaterFactory will create a new InstanceUpdater for a specific GameLift instance
type InstanceUpdaterFactory interface {
	Create(ctx context.Context, verbose bool, sshKey ssh.Signer, updateScript string, sshPort int32, instance *gamelift.Instance) (InstanceUpdater, error)
}

type instanceUpdaterFactory struct {
	logger          *slog.Logger
	gameLiftClient  GameLiftClient
	privateKeyPath  string
	buildZipPath    string
	updateOperation config.UpdateOperation
}

func NewInstanceUpdaterFactory(ctx context.Context, logger *slog.Logger, gameLiftClient GameLiftClient, args config.CLIArgs) InstanceUpdaterFactory {
	return &instanceUpdaterFactory{
		logger:          logger,
		gameLiftClient:  gameLiftClient,
		privateKeyPath:  args.PrivateKeyPath,
		buildZipPath:    args.BuildZipPath,
		updateOperation: args.GetUpdateOperation(),
	}
}

// Create will create a new instance updater that can be used to update a single instance in a GameLift fleet
func (i *instanceUpdaterFactory) Create(ctx context.Context, verbose bool, sshKey ssh.Signer, updateScript string, sshPort int32, instance *gamelift.Instance) (InstanceUpdater, error) {
	instanceLogger := i.logger.With(
		"instanceId", instance.InstanceId,
		"ipAddress", instance.IpAddress)

	sshEnabler, err := tools.NewSSHEnabler(instanceLogger, instance, i.gameLiftClient, sshKey.PublicKey(), sshPort)
	if err != nil {
		return nil, err
	}

	fileUploader, err := tools.NewFileUploader(instanceLogger, instance, i.privateKeyPath, i.GetFilesToUpload(updateScript), sshPort)
	if err != nil {
		return nil, err
	}

	commandRunner, err := tools.NewSSHCommandRunner(instanceLogger, updateScript, sshPort, sshKey, instance)
	if err != nil {
		return nil, err
	}

	progressTracker, err := NewInstanceProgressWriter(instance, verbose)
	if err != nil {
		return nil, err
	}

	return &instanceUpdater{
		sshEnabler:      sshEnabler,
		fileUploader:    fileUploader,
		commandRunner:   commandRunner,
		logger:          instanceLogger,
		progressTracker: progressTracker,
	}, nil
}

func (i *instanceUpdaterFactory) GetFilesToUpload(updateScript string) []string {
	result := make([]string, 1, 2)
	result[0] = updateScript
	if i.updateOperation == config.UpdateOperationReplaceBuild {
		result = append(result, i.buildZipPath)
	}
	return result
}
