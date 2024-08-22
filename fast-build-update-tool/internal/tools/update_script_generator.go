package tools

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
)

// InstanceUpdateScriptGenerator is used to generate a local script file which can be uploaded and run on each instance in a GameLift fleet
// The UpdateOperation provided will determine the contents of the script that is generated.
type InstanceUpdateScriptGenerator struct {
	tempBuildFile *os.File

	updateOperation   config.UpdateOperation
	localBuildZipPath string
	lockName          string
}

// NewInstanceUpdateScriptGenerator build a new InstanceUpdateScriptGenerator
func NewInstanceUpdateScriptGenerator(updateOperation config.UpdateOperation, localBuildZipPath, lockName string) *InstanceUpdateScriptGenerator {
	return &InstanceUpdateScriptGenerator{
		updateOperation:   updateOperation,
		localBuildZipPath: localBuildZipPath,
		lockName:          lockName,
	}
}

// GenerateScript will generate a script for the provided OperatingSystem.
// This function requires a slice of all of the executables that are used to run a GameServer in this specific fleet.
// The string value returned is the path on the local filesytem to the update script.
func (i *InstanceUpdateScriptGenerator) GenerateScript(ctx context.Context, operatingSystem config.OperatingSystem, executableNames []string) (filname string, err error) {
	// Actually create a file to write the update script contents to
	i.tempBuildFile, err = os.CreateTemp("", "*"+string(config.UpdateScriptForOperatingSystem(operatingSystem)))
	if err != nil {
		return "", fmt.Errorf("error creating temporary file for server update script %w", err)
	}
	defer i.tempBuildFile.Close()

	// Generate the update script
	switch operatingSystem {
	case config.OperatingSystemLinux:
		err = generateLinuxUpdateScript(i.tempBuildFile, executableNames, i.localBuildZipPath, i.lockName, i.updateOperation)
		if err != nil {
			return "", fmt.Errorf("error generating server update script %w", err)
		}

	case config.OperatingSystemWindows:
		err = generateWindowsUpdateScript(i.tempBuildFile, executableNames, i.localBuildZipPath, i.lockName, i.updateOperation)
		if err != nil {
			return "", fmt.Errorf("error generating server update script %w", err)
		}

	default:
		return "", config.UnknownOperatingSystemError(fmt.Sprint(operatingSystem))
	}

	// Return the filepath
	return i.tempBuildFile.Name(), nil
}

// Cleanup will remove the update script file generated, and clean up anything else set up by InstanceUpdateScriptGenerator
func (i *InstanceUpdateScriptGenerator) Cleanup() error {
	if i.tempBuildFile != nil {
		return os.Remove(i.tempBuildFile.Name())
	}
	return nil
}

func csvify(in []string) string {
	return strings.Join(in, ",")
}

func getIsReplaceBuildTemplateValue(updateOperation config.UpdateOperation) string {
	isReplaceBuild := ""
	if updateOperation == config.UpdateOperationReplaceBuild {
		isReplaceBuild = "replace"
	}
	return isReplaceBuild
}
