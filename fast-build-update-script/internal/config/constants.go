package config

import "log/slog"

const (
	// DefaultPortWindows is the default SSH port for Windows instances
	DefaultPortWindows int32 = 1026

	// DefaultPortLinux is the default SSH port for Linux instances
	DefaultPortLinux int32 = 22

	// AppName is the name of this application
	AppName = "fast-build-update-script"
)

// OperatingSystem is an enum of all possible GameLift operating system types
type OperatingSystem uint

// String returns a friendly version of the OperatingSystem
func (o OperatingSystem) String() string {
	switch o {
	case OperatingSystemLinux:
		return "linux"
	case OperatingSystemWindows:
		return "windows"
	default:
		return "unknown"
	}
}

const (
	// OperatingSystemUnknown is a placeholder when we don't know the OS
	OperatingSystemUnknown OperatingSystem = iota

	// OperatingSystemLinux all possible Linux instance operating systems
	OperatingSystemLinux OperatingSystem = iota

	//OperatingSystemWindows all possible Windows instance operating systems
	OperatingSystemWindows OperatingSystem = iota
)

// RemoteUser is an enum of usernames that can be used to remotely access a GameLift instance
type RemoteUser string

const (
	// RemoteUserGlUserRemote remote instance user
	RemoteUserGlUserRemote RemoteUser = "gl-user-remote"

	// RemoteUserGlUserServer remote instance user
	RemoteUserGlUserServer RemoteUser = "gl-user-server"
)

// UpdateScript is the name of the file uses to update remote instances
type UpdateScript string

const (
	// UpdateScriptWindowsName is the Windows update script file name
	UpdateScriptWindowsName UpdateScript = "update-instance.ps1"

	// UpdateScriptLinuxName is the Linux update script file name
	UpdateScriptLinuxName UpdateScript = "update-instance.sh"
)

// RemoteUploadDirectory is the directory any files will be uploaded to on the remote instance
type RemoteUploadDirectory string

const (
	// UploadDirectoryWindows the directory to upload files on a Windows instance
	UploadDirectoryWindows RemoteUploadDirectory = "C:\\Users\\gl-user-server\\"

	// UploadDirectoryLinux the directory to upload files on a Linux instance
	UploadDirectoryLinux RemoteUploadDirectory = "/tmp/"
)

// UpdateOperation is the possible update operations supported by this application
type UpdateOperation uint

const (
	// UpdateOperationReplaceBuild stop all server processes, copy a build in-place, restart server processes
	UpdateOperationReplaceBuild UpdateOperation = iota

	// UpdateOperationReplaceBuild restart all server processes
	UpdateOperationRestartProcess UpdateOperation = iota
)

// RemoteUserForOperatingSystem look up the default RemoteUser this application uses for the provided OS.
func RemoteUserForOperatingSystem(os OperatingSystem) RemoteUser {
	switch os {
	case OperatingSystemWindows:
		return RemoteUserGlUserServer
	case OperatingSystemLinux:
		return RemoteUserGlUserRemote
	default:
		slog.Warn("unknown os when looking up remote user, using default", "os", os)
		return RemoteUserGlUserRemote
	}
}

// UpdateScriptForOperatingSystem look up the update script for the provided OS.
func UpdateScriptForOperatingSystem(os OperatingSystem) UpdateScript {
	switch os {
	case OperatingSystemWindows:
		return UpdateScriptWindowsName
	case OperatingSystemLinux:
		return UpdateScriptLinuxName
	default:
		slog.Warn("unknown os when looking up update script, using default", "os", os)
		return UpdateScriptLinuxName
	}
}

// RemoteUploadDirectoryForOperatingSystem look up the remote upload directory for the provided OS.
func RemoteUploadDirectoryForOperatingSystem(os OperatingSystem) RemoteUploadDirectory {
	switch os {
	case OperatingSystemWindows:
		return UploadDirectoryWindows
	case OperatingSystemLinux:
		return UploadDirectoryLinux
	default:
		slog.Warn("unknown os when looking up remote upload directory, using default", "os", os)
		return UploadDirectoryLinux
	}
}
