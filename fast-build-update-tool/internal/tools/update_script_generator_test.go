package tools

import (
	"context"
	"os"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateLinuxReplaceBuildScript(t *testing.T) {
	updater := NewInstanceUpdateScriptGenerator(config.UpdateOperationReplaceBuild, "myarchive.zip", "lockfile")
	defer func() {
		err := updater.Cleanup()
		assert.Nil(t, err)
	}()

	filename, err := updater.GenerateScript(context.Background(), config.OperatingSystemLinux, []string{"/local/game/my-game", "/local/game/another-exe"})
	assert.Nil(t, err)

	fileBytes, err := os.ReadFile(filename)
	assert.Nil(t, err)

	fileContents := string(fileBytes)
	assert.Contains(t, fileContents, "#!/bin/bash")
	assert.Contains(t, fileContents, "ARCHIVE_NAME=myarchive.zip")
	assert.Contains(t, fileContents, `LOCKFILE="/tmp/lockfile.lock"`)
	assert.Contains(t, fileContents, "EXE_PATHS=/local/game/my-game,/local/game/another-exe")
	assert.Contains(t, fileContents, "sudo unzip -o /tmp/$ARCHIVE_NAME")
	assert.Contains(t, fileContents, "sudo rm -f $EXE_PATH;")
	assert.Contains(t, fileContents, "sudo pkill -c -f \"sudo -H -E -u gl-user-server $EXE_PATH\"")
}

func TestGenerateLinuxRestartProcessScript(t *testing.T) {
	updater := NewInstanceUpdateScriptGenerator(config.UpdateOperationRestartProcess, "", "")
	defer func() {
		err := updater.Cleanup()
		assert.Nil(t, err)
	}()

	filename, err := updater.GenerateScript(context.Background(), config.OperatingSystemLinux, []string{"/local/game/my-game"})
	assert.Nil(t, err)

	fileBytes, err := os.ReadFile(filename)
	assert.Nil(t, err)

	// make sure we don't unzip or remove anything :)
	fileContents := string(fileBytes)
	assert.Contains(t, fileContents, "#!/bin/bash")
	assert.NotContains(t, fileContents, "sudo unzip -o /tmp/$ARCHIVE_NAME")
	assert.NotContains(t, fileContents, "sudo rm -f $EXE_PATH;")
	assert.Contains(t, fileContents, "sudo pkill -c -f \"sudo -H -E -u gl-user-server $EXE_PATH\"")
}

func TestGenerateWindowsReplaceBuildScript(t *testing.T) {
	updater := NewInstanceUpdateScriptGenerator(config.UpdateOperationReplaceBuild, "myarchive.zip", "lockfile")
	defer func() {
		err := updater.Cleanup()
		assert.Nil(t, err)
	}()

	filename, err := updater.GenerateScript(context.Background(), config.OperatingSystemWindows, []string{"C:\\Game\\MyGame.exe", "C:\\Game\\other-process.exe"})
	assert.Nil(t, err)

	fileBytes, err := os.ReadFile(filename)
	assert.Nil(t, err)

	fileContents := string(fileBytes)
	assert.NotContains(t, fileContents, "#!/bin/bash")
	assert.Contains(t, fileContents, `$executablePaths="C:\Game\MyGame.exe,C:\Game\other-process.exe"`)
	assert.Contains(t, fileContents, `New-Object System.Threading.Mutex($true, "Global\lockfile", [ref]$wasLockCreated);`)
	assert.Contains(t, fileContents, `$processNames="MyGame,other-process"`)
	assert.Contains(t, fileContents, `$zipFileName="myarchive.zip"`)
	assert.Contains(t, fileContents, "Expand-Archive -Path $archivePath")
	assert.Contains(t, fileContents, "KillAll-ServerProcess $processName;")
}

func TestGenerateWindowsRestartProcessScript(t *testing.T) {
	updater := NewInstanceUpdateScriptGenerator(config.UpdateOperationRestartProcess, "myarchive.zip", "")
	defer func() {
		err := updater.Cleanup()
		assert.Nil(t, err)
	}()

	filename, err := updater.GenerateScript(context.Background(), config.OperatingSystemWindows, []string{"C:\\Game\\MyGame.exe"})
	assert.Nil(t, err)

	fileBytes, err := os.ReadFile(filename)
	assert.Nil(t, err)

	fileContents := string(fileBytes)
	assert.NotContains(t, fileContents, "#!/bin/bash")
	assert.NotContains(t, fileContents, "Expand-Archive -Path $archivePath -DestinationPath \"C:\\GameNew\\\" -Force;")
	assert.NotContains(t, fileContents, "Remove-Item -Path C:\\Game\\ -Force -Recurse;")
	assert.Contains(t, fileContents, "KillAll-ServerProcess $processName;")
}
