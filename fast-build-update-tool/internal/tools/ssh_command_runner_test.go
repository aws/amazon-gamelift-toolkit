package tools

import (
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
	"github.com/stretchr/testify/assert"
)

// TestNewSSHCommandRunnerWindows ensures we set up an ssh runner with the proper execution commands for Windows
func TestNewSSHCommandRunnerWindows(t *testing.T) {
	instance := &gamelift.Instance{OperatingSystem: config.OperatingSystemWindows}
	localUpdateScriptPath := `C:\temporary-directory\update-script.ps1`

	cmd, err := NewSSHCommandRunner(NewTestLogger(), localUpdateScriptPath, 22, nil, instance)

	assert.Nil(t, err)
	assert.Equal(t, "powershell.exe -ExecutionPolicy Bypass -File C:\\Users\\gl-user-server\\update-script.ps1", cmd.updateScriptCommand)
}

// TestNewSSHCommandRunnerWindows ensures we set up an ssh runner with the proper execution commands for Linux
func TestNewSSHCommandRunnerLinux(t *testing.T) {
	instance := &gamelift.Instance{OperatingSystem: config.OperatingSystemLinux}
	localUpdateScriptPath := `/user/local/tmp/my-script.sh`

	cmd, err := NewSSHCommandRunner(NewTestLogger(), localUpdateScriptPath, 22, nil, instance)

	assert.Nil(t, err)
	assert.Equal(t, "chmod +x /tmp/my-script.sh && /tmp/my-script.sh", cmd.updateScriptCommand)
}

// TestNewSSHCommandRunnerUnknownOS ensures we return an error when the operating system is unknown
func TestNewSSHCommandRunnerUnknownOS(t *testing.T) {
	instance := &gamelift.Instance{}
	localUpdateScriptPath := `/user/local/tmp/my-script.sh`

	_, err := NewSSHCommandRunner(NewTestLogger(), localUpdateScriptPath, 22, nil, instance)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "argument operatingSystem was invalid")
}
