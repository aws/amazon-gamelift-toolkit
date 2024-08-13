package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/gamelift"
	"github.com/stretchr/testify/assert"
)

var (
	currentDir   string
	buildZipPath string
)

func init() {
	currentDir, _ = os.Getwd()
	buildZipPath = filepath.Join(currentDir, "testdata", "game-executable.zip")
}

func TestValidateZipWindowsValidZip(t *testing.T) {
	z := NewZipValidator(buildZipPath)
	fleet := &gamelift.Fleet{OperatingSystem: config.OperatingSystemWindows, ExecutablePaths: []string{"C:\\game\\bin\\server.exe"}}

	err := z.ValidateZip(context.Background(), fleet)

	assert.Nil(t, err)
}

func TestValidateZipLinuxValidZip(t *testing.T) {
	z := NewZipValidator(buildZipPath)
	fleet := &gamelift.Fleet{OperatingSystem: config.OperatingSystemLinux, ExecutablePaths: []string{"/local/game/bin/server.exe"}}

	err := z.ValidateZip(context.Background(), fleet)

	assert.Nil(t, err)
}

func TestValidateZipInvalidZipNoFile(t *testing.T) {
	z := NewZipValidator("file does not exist")
	fleet := &gamelift.Fleet{OperatingSystem: config.OperatingSystemLinux, ExecutablePaths: []string{"/local/game/bin/server.exe"}}

	err := z.ValidateZip(context.Background(), fleet)

	assert.ErrorContains(t, err, "error opening zip file")
}

func TestValidateZipInvalidZipNoExecutable(t *testing.T) {
	z := NewZipValidator(buildZipPath)
	fleet := &gamelift.Fleet{OperatingSystem: config.OperatingSystemLinux, ExecutablePaths: []string{"/local/game/different-server.exe"}}

	err := z.ValidateZip(context.Background(), fleet)

	assert.ErrorContains(t, err, "zip file does not contain executable")
}
