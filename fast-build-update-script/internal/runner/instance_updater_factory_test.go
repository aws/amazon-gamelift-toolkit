package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/gamelift"
	"github.com/stretchr/testify/assert"
)

func TestGetFilesToUploadRestartProcess(t *testing.T) {
	zipPath := "myfile.zip"
	updateScript := "update-script.sh"

	i := &instanceUpdaterFactory{updateOperation: config.UpdateOperationRestartProcess, buildZipPath: zipPath}

	filesToUpload := i.GetFilesToUpload(updateScript)

	assert.Len(t, filesToUpload, 1)
	assert.Equal(t, updateScript, filesToUpload[0])
}

func TestGetFilesToUploadReplaceBuild(t *testing.T) {
	zipPath := "myfile.zip"
	updateScript := "update-script.sh"

	i := &instanceUpdaterFactory{updateOperation: config.UpdateOperationReplaceBuild, buildZipPath: zipPath}

	filesToUpload := i.GetFilesToUpload(updateScript)

	assert.Len(t, filesToUpload, 2)
	assert.Equal(t, updateScript, filesToUpload[0])
	assert.Equal(t, zipPath, filesToUpload[1])
}

func TestCreate(t *testing.T) {
	signer, privateKeyPath := generatePrivateSSHKey()
	defer os.Remove(privateKeyPath)

	currentDir, err := os.Getwd()
	assert.Nil(t, err)

	buildZipPath := filepath.Join(currentDir, "testdata", "game-executable.zip")

	factory := NewInstanceUpdaterFactory(context.Background(), NewTestLogger(), &GameLiftClientMock{}, config.CLIArgs{
		FleetId:        fleetId,
		IpRange:        "0.0.0.0/0",
		BuildZipPath:   buildZipPath,
		SSHPort:        22,
		InstanceIds:    make([]string, 0),
		RestartProcess: false,
		LockName:       "test",
		Verbose:        false,
		PrivateKeyPath: privateKeyPath,
	})

	instanceUpdater, err := factory.Create(context.Background(), true, signer, "update-script", 22, &gamelift.Instance{OperatingSystem: config.OperatingSystemLinux})
	assert.Nil(t, err)
	assert.NotNil(t, instanceUpdater)
}
