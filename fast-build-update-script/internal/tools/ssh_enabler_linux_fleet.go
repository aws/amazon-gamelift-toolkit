package tools

import (
	"fmt"
	"regexp"
)

var (
	linuxNewCommandRegex = regexp.MustCompile(`(?m)sh-\d\.\d\$`)
)

func IsNewCommandOutputLinux(output string) bool {
	return linuxNewCommandRegex.MatchString(output)
}

func linuxSSHEnableCommands(localPublicKey string) []string {
	return []string{
		"sudo touch /home/gl-user-remote/.ssh/authorized_keys;\n",
		fmt.Sprintf("echo \"%s\" | sudo tee /home/gl-user-remote/.ssh/authorized_keys;\n", localPublicKey),
		"cat /etc/ssh/ssh_host_ed25519_key.pub;\n",
		"exit;\n",
	}
}
