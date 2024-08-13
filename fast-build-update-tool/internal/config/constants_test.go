package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperatingSystemString(t *testing.T) {
	assert.Equal(t, "linux", OperatingSystemLinux.String())
	assert.Equal(t, "windows", OperatingSystemWindows.String())
	assert.Equal(t, "unknown", OperatingSystemUnknown.String())
}

func TestRemoteUserForOperatingSystem(t *testing.T) {
	assert.Equal(t, "gl-user-server", string(RemoteUserForOperatingSystem(OperatingSystemWindows)))

	assert.Equal(t, "gl-user-remote", string(RemoteUserForOperatingSystem(OperatingSystemLinux)))

	assert.Equal(t, "gl-user-remote", string(RemoteUserForOperatingSystem(OperatingSystemUnknown)))
}

func TestUpdateScriptForOperatingSystem(t *testing.T) {
	assert.Equal(t, "update-instance.ps1", string(UpdateScriptForOperatingSystem(OperatingSystemWindows)))

	assert.Equal(t, "update-instance.sh", string(UpdateScriptForOperatingSystem(OperatingSystemLinux)))

	assert.Equal(t, "update-instance.sh", string(UpdateScriptForOperatingSystem(OperatingSystemUnknown)))
}

func TestRemoteUploadDirectoryForOperatingSystem(t *testing.T) {
	assert.Equal(t, "C:\\Users\\gl-user-server\\", string(RemoteUploadDirectoryForOperatingSystem(OperatingSystemWindows)))

	assert.Equal(t, "/tmp/", string(RemoteUploadDirectoryForOperatingSystem(OperatingSystemLinux)))

	assert.Equal(t, "/tmp/", string(RemoteUploadDirectoryForOperatingSystem(OperatingSystemUnknown)))
}
