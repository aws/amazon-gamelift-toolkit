// config holds any logic around application configuration, logging, common constants, etc...
package config

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// CLIArgs holds the parsed and validated args the user passed to the application
type CLIArgs struct {
	// FleetId is the id of the fleet the user would like to update
	FleetId string
	// IpRange is the range of IP address that are allowed remote access to the GameLift fleet
	IpRange string
	// BuildZipPath is the path on the local filesystem to the build zip file
	BuildZipPath string
	// PrivateKeyPath is the path on the local filesystem to the private SSH key that will be used to interact with remote instances
	PrivateKeyPath string
	// SSHPort is the port that will be opened for SSH use on any remote instances
	SSHPort int
	// InstanceIds is an optional allow list of instance ids to update in GameLift
	InstanceIds []string
	// RestartProcess optional flag to skip uploading and replacing a build, and simply restart the server process on remote instances
	RestartProcess bool
	// LockName is an optional override to change the name of the lock file used on remote servers in-case of deadlock.
	LockName string
	// Verbose is an optional argument to provide more verbose application logs
	Verbose bool

	instanceIdsRaw string
}

// ParseAndValidateCLIArgs will parse the input slice of string arguments, and validate them
func ParseAndValidateCLIArgs(cliArgs []string) (CLIArgs, error) {
	result, err := ParseArgs(cliArgs)
	if err != nil {
		return result, err
	}

	return result, result.Validate()
}

const (
	argFleetId        = "fleet-id"
	argIpRange        = "ip-range"
	argBuildZipPath   = "zip-path"
	argPrivateKey     = "private-key"
	argSSHPort        = "ssh-port"
	argInstanceIds    = "instance-ids"
	argRestartProcess = "restart-process"
	argLockName       = "lock-name"
	argVerbose        = "verbose"
)

// ParseArgs will parse the input slice of string arguments into CLIArgs
func ParseArgs(args []string) (CLIArgs, error) {
	result := CLIArgs{}

	flags := flag.NewFlagSet(AppName, flag.ContinueOnError)

	// Define required arguments
	flags.StringVar(&result.FleetId, argFleetId, "", "[Required] The ID of the GameLift Fleet to update")
	flags.StringVar(&result.IpRange, argIpRange, "", "[Required] Your local IP Address, needed to open ports on the fleet for remote connections (eg. 127.0.0.1/32)")
	flags.StringVar(&result.BuildZipPath, argBuildZipPath, "", "[Required] The path to the zip file containing your build")
	flags.StringVar(&result.PrivateKeyPath, argPrivateKey, "", "[Required] The local path to a private key to be used with SSH")

	// Define optional arguments
	flags.IntVar(&result.SSHPort, argSSHPort, 0, "[Optional] The port to open for SSH on the fleet. This option is for Windows remote instances only. It will default to 1026.")
	flags.StringVar(&result.instanceIdsRaw, argInstanceIds, "", "[Optional] A list of instance ids to update separated by comma. If not provided all instances will be updated")
	flags.BoolVar(&result.RestartProcess, argRestartProcess, false, "[Optional] Flag to restart existing game server processes on a server, and skip uploading a new build and replacing the old build.")
	flags.StringVar(&result.LockName, argLockName, AppName, "[Optional] This should only be set if you encounter a deadlock. This should not be set in typical application use. Set this argument to manually override the lock file name used on the server if your application gets stuck in an update deadlock.")
	flags.BoolVar(&result.Verbose, argVerbose, false, "[Optional] Write more verbose logs as output")

	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s --%s FLEET_ID --%s IP_RANGE --%s BUILD_ZIP_PATH --%s PRIVATE_KEY \n", os.Args[0], argFleetId, argIpRange, argBuildZipPath, argPrivateKey)
		flags.PrintDefaults()
	}

	// If nothing was passed at all, show the usage instructions
	if len(args) <= 1 {
		flags.Usage()
		return result, flag.ErrHelp
	}

	// Parse the arguments (without the application exe in the slice)
	err := flags.Parse(args[1:])
	if err != nil {
		return result, err
	}

	// Split instance id CSV into a slice if provided
	if result.instanceIdsRaw != "" {
		result.InstanceIds = strings.Split(result.instanceIdsRaw, ",")
	}

	return result, nil
}

// Validate that all of the CLIArgs are valid
func (c *CLIArgs) Validate() (err error) {
	if c.FleetId == "" {
		err = errors.Join(err, missingArgumentError(argFleetId))
	}

	if c.IpRange == "" {
		err = errors.Join(err, missingArgumentError(argIpRange))

	} else if !isValidIpRange(c.IpRange) {
		err = errors.Join(err, invalidArgumentError(argIpRange, "must be a valid IP range"))
	}

	// We do not need to validate a build zip file if we are just restarting the process
	if c.RestartProcess {
		if c.BuildZipPath != "" {
			err = errors.Join(err, invalidArgumentError(argBuildZipPath, "zip file provided along with restart process flag"))
		}
	} else {
		if c.BuildZipPath == "" {
			err = errors.Join(err, missingArgumentError(argBuildZipPath))

		} else if !doesFileExist(c.BuildZipPath) {
			err = errors.Join(err, missingFileError(argBuildZipPath))

		}
	}

	if c.PrivateKeyPath == "" {
		err = errors.Join(err, missingArgumentError(argPrivateKey))

	} else if !doesFileExist(c.PrivateKeyPath) {
		err = errors.Join(err, missingFileError(argPrivateKey))
	}

	return err
}

// GetUpdateOperation will return what update operation the CLIArgs have instructed the app to take
func (c *CLIArgs) GetUpdateOperation() UpdateOperation {
	// Currently only two options here (restart processes, or replace the build on all instances)
	if c.RestartProcess {
		return UpdateOperationRestartProcess
	}
	return UpdateOperationReplaceBuild
}

func doesFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isValidIpRange(ipRange string) bool {
	_, _, cidrErr := net.ParseCIDR(ipRange)
	return cidrErr == nil
}
