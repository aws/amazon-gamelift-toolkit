package runner

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/tools"
	"golang.org/x/crypto/ssh"
)

var UpdateFailedError error = errors.New("failed to update one or more instances")

// FleetUpdater coordinates applying updates to instances in a specific GameLift fleet
type FleetUpdater struct {
	args config.CLIArgs

	logger *slog.Logger

	gameLiftClient         GameLiftClient
	updateScriptGenerator  *tools.InstanceUpdateScriptGenerator
	sshConfigManager       *tools.SSHConfigManager
	zipValidator           *tools.ZipValidator
	instanceUpdaterFactory InstanceUpdaterFactory
	reportWriter           *FleetUpdateReportWriter
}

// NewFleetUpdater will build a new FleetUpdater using command line arguments
func NewFleetUpdater(ctx context.Context, logger *config.ApplicationLogger, args config.CLIArgs) (*FleetUpdater, error) {
	slogger := logger.Logger.With("fleetId", args.FleetId)

	gameLift, err := gamelift.NewGameLiftClient(ctx, logger.AwsLogger)
	if err != nil {
		return nil, err
	}

	return &FleetUpdater{
		args:                   args,
		gameLiftClient:         gameLift,
		logger:                 slogger,
		updateScriptGenerator:  tools.NewInstanceUpdateScriptGenerator(args.GetUpdateOperation(), args.BuildZipPath, args.LockName),
		sshConfigManager:       tools.NewSSHConfigManager(slogger, args.PrivateKeyPath, args.SSHPort),
		zipValidator:           tools.NewZipValidator(args.BuildZipPath),
		instanceUpdaterFactory: NewInstanceUpdaterFactory(ctx, slogger, gameLift, args),
		reportWriter:           NewFleetUpdateReportWriter(args.FleetId, args.Verbose),
	}, nil
}

// UpdateInstances will perform any actions necessary to update instances in a GameLift fleet
func (f *FleetUpdater) UpdateInstances(ctx context.Context) (*FleetUpdateResults, error) {
	f.logger.Info("starting fleet update process")

	f.reportWriter.Preparing()

	fleet, err := f.lookupFleet(ctx)
	if err != nil {
		return nil, err
	}

	err = f.validateZipFile(ctx, fleet)
	if err != nil {
		return nil, err
	}

	sshPort, err := f.ensureSSHPortIsSet(ctx, fleet.OperatingSystem)
	if err != nil {
		return nil, err
	}

	err = f.ensureSSHPortIsOpenForFleet(ctx, sshPort)
	if err != nil {
		return nil, err
	}

	updateScript, err := f.generateUpdateScript(ctx, fleet)
	if err != nil {
		return nil, err
	}

	sshKey, err := f.loadSSHKey(ctx)
	if err != nil {
		return nil, err
	}

	instances, err := f.getInstances(ctx)
	if err != nil {
		return nil, err
	}

	return f.updateInstances(ctx, instances, sshKey, sshPort, fleet.OperatingSystem, updateScript)
}

// lookupFleet will verify the fleet exists, and fetch any relevant data we need to perform an update
func (f *FleetUpdater) lookupFleet(ctx context.Context) (*gamelift.Fleet, error) {
	fleet, err := f.gameLiftClient.GetFleet(ctx, f.args.FleetId)
	if err != nil {
		return fleet, fmt.Errorf("error looking up fleet: %w", err)
	}

	f.logger.Debug("looking up fleet attributes", "os", fleet.OperatingSystem, "executables", fleet.ExecutablePaths)

	return fleet, nil
}

// validateZipFile will validate that the zip file provided by the user is valid for the given fleet
func (f *FleetUpdater) validateZipFile(ctx context.Context, fleet *gamelift.Fleet) error {
	// If the user is restarting server processes, we don't have a zip file to validate
	if f.args.RestartProcess {
		f.logger.Debug("running as a restart process update, skipping zip file validation")
		return nil
	}

	err := f.zipValidator.ValidateZip(ctx, fleet)
	if err != nil {
		return fmt.Errorf("error validating zip file: %w", err)
	}

	f.logger.Debug("done validating zip file")

	return nil
}

// ensureSSHPortIsSet will make sure we have a valid SSH port set for the operating system running in this fleet
func (f *FleetUpdater) ensureSSHPortIsSet(ctx context.Context, os config.OperatingSystem) (int32, error) {
	port, err := f.sshConfigManager.DeterminePort(os)
	if err != nil {
		return port, fmt.Errorf("error determining ssh port %w", err)
	}

	f.logger.Debug("done determining SSH port", "port", port)

	return port, nil
}

// ensureSSHPortIsOpenForFleet will update GameLift configuration to verify the ssh port is open for the IP range provided by the user
func (f *FleetUpdater) ensureSSHPortIsOpenForFleet(ctx context.Context, sshPort int32) error {
	err := f.gameLiftClient.OpenPortForFleet(ctx, f.args.FleetId, sshPort, f.args.IpRange)
	if err != nil {
		return fmt.Errorf("error opening port for fleet %w", err)
	}

	f.logger.Debug("done ensuring SSH port is open for fleet")

	return nil
}

// loadSSHKey will load the SSH key provided by the user
func (f *FleetUpdater) loadSSHKey(ctx context.Context) (ssh.Signer, error) {
	signer, err := f.sshConfigManager.LoadKey(ctx)
	if err != nil {
		return signer, fmt.Errorf("error loading private ssh key %w", err)
	}

	f.logger.Debug("done loading ssh key")

	return signer, nil
}

// generateUpdateScript will generate the script we use to update individual instances in the fleet.
// This script will be uploaded and run on each individual instance later in this process.
func (f *FleetUpdater) generateUpdateScript(ctx context.Context, fleet *gamelift.Fleet) (string, error) {
	scriptPath, err := f.updateScriptGenerator.GenerateScript(ctx, fleet.OperatingSystem, fleet.ExecutablePaths)
	if err != nil {
		return scriptPath, fmt.Errorf("error generating update script: %w", err)
	}

	f.logger.Debug("done generating update script", "scriptPath", scriptPath)

	return scriptPath, nil
}

// getInstances will load any relevant instances for this update operation
func (f *FleetUpdater) getInstances(ctx context.Context) ([]*gamelift.Instance, error) {
	instances, err := f.gameLiftClient.GetInstances(ctx, f.args.FleetId, f.args.InstanceIds)
	if err != nil {
		return instances, fmt.Errorf("error fetching instances for fleet: %w", err)
	}

	f.logger.Debug("done loading instances in GameLift fleet", "instanceCount", len(instances))

	return instances, nil
}

// updateInstances will actually run through the process of updating each instance in the fleet
func (f *FleetUpdater) updateInstances(ctx context.Context, instances []*gamelift.Instance, sshKey ssh.Signer, sshPort int32, os config.OperatingSystem, updateScript string) (*FleetUpdateResults, error) {
	f.logger.Debug("updating instances in GameLift fleet")

	f.reportWriter.StartUpdatingInstances(len(instances))

	results := &FleetUpdateResults{
		InstancesFound:        len(instances),
		InstancesFailedUpdate: make([]string, 0, len(instances)),
	}

	for _, instance := range instances {
		err := f.updateInstance(ctx, sshKey, sshPort, updateScript, instance)
		if err != nil {
			// If we fail to update an instance, log the error and continue. We may still be able to update other instances in the fleet
			slog.Error("Error updating remote instance", "error", err, "instanceId", instance.InstanceId)
			results.InstancesFailedUpdate = append(results.InstancesFailedUpdate, instance.InstanceId)
			continue
		}

		results.InstancesUpdated = results.InstancesUpdated + 1
	}

	// We're done updating instances, write the report out for the user
	f.reportWriter.ReportResults(results)

	// If any instances failed to update, ensure that we return an error
	if len(results.InstancesFailedUpdate) > 0 {
		return results, UpdateFailedError
	}

	return results, nil
}

// updateInstance update an individual instance in the fleet
func (f *FleetUpdater) updateInstance(ctx context.Context, sshKey ssh.Signer, sshPort int32, updateScript string, instance *gamelift.Instance) error {
	// create a fleet updater
	instanceUpdater, err := f.instanceUpdaterFactory.Create(ctx, f.args.Verbose, sshKey, updateScript, sshPort, instance)
	if err != nil {
		return fmt.Errorf("error setting up instance updater: %w", err)
	}

	// update the instance
	err = instanceUpdater.Update(ctx)
	if err != nil {
		return fmt.Errorf("error updating instance: %w", err)
	}

	return nil
}

func (f *FleetUpdater) Cleanup() {
	f.logger.Debug("cleaning up fleet updater resources")

	if f.updateScriptGenerator != nil {
		err := f.updateScriptGenerator.Cleanup()
		if err != nil {
			f.logger.Warn("error cleaning up local update script", "err", err)
		}
	}

	f.logger.Debug("done cleaning up fleet updater resources")
}
