package tools

import "os/exec"

// verifyExe will verify that the user has the provided executable in their path
func verifyExe(exePath string) error {
	_, err := exec.LookPath(scpCommand)
	return err
}
