package tools

import (
	"fmt"
	"regexp"
)

var (
	windowsNewCommandRegex = regexp.MustCompile(`(?m)PS\s.*\>`)
)

func IsNewCommandOutputWindows(output string) bool {
	return windowsNewCommandRegex.MatchString(output)
}

func windowsSSHEnableCommands(localPublicKey string, sshPort int32) []string {
	return []string{
		fmt.Sprintf("$port=\"%d\";\r\n", sshPort),
		fmt.Sprintf("$publicKey=\"%s\";\r\n", localPublicKey),
		windowsInstallSSHPowershellScript,
		"Get-Content -Path C:\\ProgramData\\ssh\\ssh_host_ed25519_key.pub;\r\n",
		"exit;\r\n",
	}
}

const windowsInstallSSHPowershellScript = `
$isSSHRunning = net start | Select-String -Pattern OpenSSH;
if (!$isSSHRunning) {
	Write-Host "Setting up OpenSSH"

	New-NetFirewallRule -Name sshd -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort $port -ErrorAction SilentlyContinue;
	
	if (!(Test-Path "C:\OpenSSH-Win64.zip")) {
		[Net.ServicePointManager]::SecurityProtocol = "tls12, tls11, tls";
		Invoke-WebRequest -Uri "https://github.com/PowerShell/Win32-OpenSSH/releases/download/v8.9.1.0p1-Beta/OpenSSH-Win64.zip" -OutFile "C:\OpenSSH-Win64.zip";
	}

	if (!(Test-Path "C:\Program Files\OpenSSH-Win64")) {
		Expand-Archive C:\OpenSSH-Win64.zip -DestinationPath "C:\Program Files";
		Remove-Item -Path "C:\OpenSSH-Win64.zip";
	}
		
	Set-Location -Path "C:\Program Files\OpenSSH-Win64\";
	powershell.exe -ExecutionPolicy Bypass -File install-sshd.ps1;
	
	net start sshd;
}


$sshConfig = "C:\ProgramData\ssh\sshd_config";

if (!(Get-Content $sshConfig | Select-String -Pattern "^Port $port$")) {
	# Remove old port line from the file if there is one
	(Get-Content $sshConfig) | Where-Object { $_ -notmatch '^Port ' } | Set-Content $sshConfig

	# Append new port line to the file and restart the SSH service
	$portLine = "Port "+$port;
	$oldContent = Get-Content -Path $sshConfig;
	$newContent = $portLine, $oldContent;
	$newContent | Set-Content -Path $sshConfig;

	# Restart service
	net stop sshd;
	net start sshd;
}

# Add SSH to the path
$path = [System.Environment]::GetEnvironmentVariable("PATH") -split ";"
if (!($path -contains "C:\Program Files\OpenSSH-Win64")) {
	setx PATH "$env:PATH;C:\Program Files\OpenSSH-Win64" -m;
}

# Create C:\Users\gl-user-server\.ssh\authorized_keys file if it doesn't already exist
if (!(Test-Path "C:\Users\gl-user-server\.ssh\")) {
	New-Item -Path "C:\Users\gl-user-server\.ssh\" -ItemType Directory;
}

if (!(Test-Path "C:\Users\gl-user-server\.ssh\authorized_keys")) {
	New-Item -Path "C:\Users\gl-user-server\.ssh\authorized_keys" -ItemType File;
}

# Add public key to authorized keys file
if (!(Select-String -Path C:\Users\gl-user-server\.ssh\authorized_keys -Pattern $publicKey)) {
	Add-Content -Path "C:\Users\gl-user-server\.ssh\authorized_keys" -Value $publicKey;
}
`
