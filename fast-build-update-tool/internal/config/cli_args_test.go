package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	currentDir     string
	buildZipPath   string
	privateKeyPath string
)

func init() {
	currentDir, _ = os.Getwd()
	buildZipPath = filepath.Join(currentDir, "testdata", "game-executable.zip")
	privateKeyPath = filepath.Join(currentDir, "testdata", "fake-ssh-key")
}

// TestParseAndValidateCLIArgs validates that this function properly parses valid args without error
func TestParseAndValidateCLIArgs(t *testing.T) {
	args, err := ParseAndValidateCLIArgs([]string{
		"appName.exe",
		"--fleet-id", "1234",
		"--ip-range", "0.0.0.0/0",
		"--private-key", privateKeyPath,
		"--zip-path", buildZipPath})

	assert.Nil(t, err)
	assert.Equal(t, "1234", args.FleetId)
	assert.Equal(t, "0.0.0.0/0", args.IpRange)
	assert.Equal(t, privateKeyPath, args.PrivateKeyPath)
	assert.Equal(t, buildZipPath, args.BuildZipPath)
}

// TestParseAndValidateCLIArgsArgsEmpty validates that we return the help/usage error when no args are passed
func TestParseAndValidateCLIArgsArgsEmpty(t *testing.T) {
	_, err := ParseAndValidateCLIArgs([]string{})
	assert.Equal(t, flag.ErrHelp, err)

	_, err = ParseArgs([]string{"appName.exe"})
	assert.Equal(t, flag.ErrHelp, err)
}

// TestParseAndValidationError validates that we return any validation errors when args are invalid
func TestParseAndValidationError(t *testing.T) {
	_, err := ParseAndValidateCLIArgs([]string{"appName.exe", "--fleet-id", "1234"})
	assert.NotNil(t, err)
}

// TestParseArgs validates that parse args properly parses all possible arguments
func TestParseArgs(t *testing.T) {
	fleetId := "my-fleet"
	ipRange := "0.0.0.0/0"
	zipPath := "mybuild.zip"
	privateKey := "mykey.pem"
	sshPort := 22
	instanceIds := "1,2"
	lockName := "mycustomlock"
	args, err := ParseArgs([]string{"appName.exe",
		"--fleet-id", fleetId,
		"--ip-range", ipRange,
		"--zip-path", zipPath,
		"--private-key", privateKey,
		"--ssh-port", strconv.Itoa(sshPort),
		"--instance-ids", instanceIds,
		"--restart-process",
		"--lock-name", lockName,
		"--verbose"})

	assert.Nil(t, err)
	assert.Equal(t, fleetId, args.FleetId)
	assert.Equal(t, ipRange, args.IpRange)
	assert.Equal(t, zipPath, args.BuildZipPath)
	assert.Equal(t, privateKey, args.PrivateKeyPath)
	assert.Equal(t, sshPort, args.SSHPort)
	assert.Contains(t, args.InstanceIds, "1")
	assert.Contains(t, args.InstanceIds, "2")
	assert.True(t, args.RestartProcess)
	assert.Equal(t, lockName, args.LockName)
	assert.True(t, args.Verbose)
}

// TestValidateValidArgs validates that valid args return no errors
func TestValidateValidArgs(t *testing.T) {
	args := &CLIArgs{
		FleetId:        "fleet-id",
		IpRange:        "127.0.0.1/0",
		BuildZipPath:   buildZipPath,
		PrivateKeyPath: privateKeyPath,
	}

	err := args.Validate()

	assert.Nil(t, err)
}

// TestValidateValidArgsRestartProcess validates that no error is returned when RestartProcess is true without a zip path
func TestValidateValidArgsRestartProcess(t *testing.T) {
	args := &CLIArgs{
		FleetId:        "fleet-id",
		IpRange:        "127.0.0.1/0",
		RestartProcess: true,
		PrivateKeyPath: privateKeyPath,
		// skip BuildZipPath
	}

	err := args.Validate()
	assert.Nil(t, err)
}

// TestValidateValidArgsRestartProcessWithBuildZip validates that an error occurs when a zip file is provided during a RestartProcess update.
func TestValidateValidArgsRestartProcessWithBuildZip(t *testing.T) {
	args := &CLIArgs{
		FleetId:        "fleet-id",
		IpRange:        "127.0.0.1/0",
		RestartProcess: true,
		PrivateKeyPath: privateKeyPath,
		BuildZipPath:   buildZipPath,
	}

	err := args.Validate()

	assert.Contains(t, err.Error(), "zip file provided along with restart process flag")
}

// TestValidateEmptyStruct validates that errors are returned for all required arguments
func TestValidateEmptyStruct(t *testing.T) {
	args := &CLIArgs{}

	err := args.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required argument fleet-id")
	assert.Contains(t, err.Error(), "missing required argument ip-range")
	assert.Contains(t, err.Error(), "missing required argument zip-path")
	assert.Contains(t, err.Error(), "missing required argument private-key")
}

// TestValidateIpRange validates that IP range validation works correctly
func TestValidateIpRange(t *testing.T) {
	args := &CLIArgs{IpRange: "127.0.0.1"}

	err := args.Validate()

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "argument ip-range was invalid: must be a valid IP range")
}

// TestValidateFilesDoNotExist ensures that proper errors are returned when non-existent file arguments are passed in
func TestValidateFilesDoNotExist(t *testing.T) {
	args := &CLIArgs{BuildZipPath: "not a real zip file", PrivateKeyPath: "not a real private key file"}

	err := args.Validate()

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "argument zip-path was invalid: could not find file")
	assert.ErrorContains(t, err, "argument private-key was invalid: could not find file")
}

// TestGetUpdateOperationRestartProcess validates that GetUpdateOperation returns the proper value for a restart process update
func TestGetUpdateOperationRestartProcess(t *testing.T) {
	args := &CLIArgs{RestartProcess: true}
	assert.Equal(t, args.GetUpdateOperation(), UpdateOperationRestartProcess)
}

// TestGetUpdateOperationRestartProcess validates that GetUpdateOperation returns the proper value for a replace build update
func TestGetUpdateOperationReplaceBuild(t *testing.T) {
	args := &CLIArgs{}
	assert.Equal(t, args.GetUpdateOperation(), UpdateOperationReplaceBuild)
}
