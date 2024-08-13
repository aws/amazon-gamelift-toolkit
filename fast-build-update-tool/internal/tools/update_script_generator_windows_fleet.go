package tools

import (
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
)

// generateLinuxUpdateScript is used to generate an update script for a Windows fleet
// updateOperation configures which type of update script will be generated
func generateWindowsUpdateScript(writer io.Writer, executablePaths []string, localBuildZipPath, lockName string, updateOperation config.UpdateOperation) error {
	template, err := template.New("windows-update-template").Parse(windowsUpdateScriptTemplate)
	if err != nil {
		return err
	}

	processNames := make([]string, len(executablePaths))
	for i, executablePath := range executablePaths {
		parts := strings.Split(executablePath, "\\")
		exeName := parts[len(parts)-1]
		processNames[i] = strings.Replace(exeName, ".exe", "", -1)
	}

	return template.Execute(writer, map[string]string{
		"ArchiveName":     filepath.Base(localBuildZipPath),
		"ExecutablePaths": csvify(executablePaths),
		"ProcessNames":    csvify(processNames),
		"IsReplaceBuild":  getIsReplaceBuildTemplateValue(updateOperation),
		"LockName":        lockName,
	})
}

const windowsUpdateScriptTemplate = `
$ErrorActionPreference = "Stop";    

[bool]$wasLockCreated = $false;
[System.Threading.Mutex]$mutex;

[Reflection.Assembly]::LoadWithPartialName("System.IO.Compression.ZipFile")

$baseDir="C:\Game\";
$unzipDir="C:\GameNew\";

$executablePaths="{{ .ExecutablePaths }}" -split ",";
$processNames="{{ .ProcessNames }}" -split ",";
$zipFileName="{{ .ArchiveName }}";
$archivePath="C:\Users\gl-user-server\$zipFileName";

try { 

$mutex = New-Object System.Threading.Mutex($true, "Global\{{ .LockName }}", [ref]$wasLockCreated);
if (!$wasLockCreated)
{
	Write-Host "ERROR! Couldn't acquire update lock, exiting...";
	exit 1;
}
Write-Host "Acquired update lock";

function KillAll-ServerProcess {
	param (
        [string]$ProcessToKill
    )

	Write-Host "Stopping all processes with name: $ProcessToKill";

	$serverProcesses = Get-Process -Name $ProcessToKill -ErrorAction SilentlyContinue;
	if ($serverProcesses) {
		Write-Host "Stopping the old server processes";
		foreach ($process in $serverProcesses) {
			Write-Host "Stopping the process with id: " $process.Id;

			# Stop the process
			Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue;

			# Wait for the process to exit
			Wait-Process -Id $process.Id -ErrorAction SilentlyContinue;
			
			Write-Host "Done stopping the process with id: " $process.Id;
		}
	} else {
		Write-Host "No running process found: $ProcessToKill";
	}
}

Write-Host "===========================================================";
Write-Host "Ending running server processes";
Write-Host "===========================================================";

{{if .IsReplaceBuild}}

foreach ($executablePath in $executablePaths) {
	if (Test-Path $executablePath) {
		Write-Host "Moving old executable to $executablePath-old";
		Move-Item -Force -Path $executablePath -Destination $executablePath-old;
	} else {
		Write-Host "Executable $executablePath not found";
	}
}

{{end}}

foreach ($processName in $processNames) {
	KillAll-ServerProcess $processName;
}

{{if .IsReplaceBuild}}

Start-Sleep -Seconds 5;

Write-Host "===========================================================";
Write-Host "Removing files found in the build zip from the server";
Write-Host "===========================================================";

$zip=[System.IO.Compression.ZipFile]::OpenRead($archivePath);
try {
	$zip.Entries | ForEach-Object {
		$isDirectory= $_.FullName[-1] -eq '/' -or $_.FullName[-1] -eq '\';
		if (!$isDirectory) {
			$fileName= $_.FullName -replace '/', '\';
			$removePath=$baseDir + $fileName;

			if (Test-Path $removePath)
			{
				Write-Host "Removing old build file: $removePath";
				Remove-Item -Path $removePath -Force;
			} else {
				Write-Host "File from build zip file: $removePath, not seen on the server.";
			}
		}
	}
} catch {
	Write-Host "An unexpected error occurred:"
	Write-Host $_
	throw $_
} finally {
	$zip.Dispose();
}


Write-Host "===========================================================";
Write-Host "Expanding $archivePath to $unzipDir";
Write-Host "===========================================================";
Expand-Archive -Path $archivePath -DestinationPath $unzipDir -Force;

Write-Host "===========================================================";
Write-Host "Moving files from $unzipDir to $baseDir";
Write-Host "===========================================================";

Get-ChildItem -Path $unzipDir -Recurse | ForEach-Object {
	$name=$_.FullName;
    $destination = Join-Path -Path $baseDir -ChildPath $name.Substring($unzipDir.Length);

	if ($executablePaths -contains $destination) {
		Write-Host "Found executable $destination, this must be copied last";
	} else {
		$destinationDir = Split-Path -Path $destination -Parent;
		if (Test-Path -Path $destinationDir) {
		} else {
			Write-Host "Missing host dir creating it: $destinationDir";
    	    New-Item -Path $destinationDir -ItemType Directory | Out-Null;
    	}

		if (Test-Path $destination) {
		} else {
			Write-Host "Moving file to $destination";
			Move-Item -Path $_.FullName -Destination $destination -Force;
		}
	}
}

foreach ($executablePath in $executablePaths) {
	$unzipPath = Join-Path -Path $unzipDir -ChildPath $executablePath.Substring($baseDir.Length);
	Write-Host "Moving executable file $unzipPath to $executablePath";
	Move-Item -Path $unzipPath -Destination $executablePath -Force;

	if (Test-Path $executablePath-old) {
		Write-Host "Removing $executablePath-old";
		Remove-Item -Path $executablePath-old;
	}
}

{{end}}

} catch {
	Write-Host "An unexpected error occurred:"
	Write-Host $_
	throw $_

} finally {
	if ($wasLockCreated -and $null -ne $mutex) {
		$mutex.ReleaseMutex()
		$mutex.Dispose()
		Write-Host "Update lock released"
	}

{{if .IsReplaceBuild}}
	Write-Host "Cleaning up archive $archivePath";
	Remove-Item -Path $archivePath -Force;
	if (Test-Path $unzipDir) {
		Remove-Item -Recurse -Force -Path $unzipDir;
	}
{{end}}

	Write-Host "Cleaning up update script $PSCommandPath";
	Remove-Item $PSCommandPath -Force;
}
`
