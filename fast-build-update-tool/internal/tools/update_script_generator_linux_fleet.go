package tools

import (
	"io"
	"path/filepath"
	"text/template"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
)

// generateLinuxUpdateScript is used to generate an update script for a Linux fleet
// updateOperation configures which type of update script will be generated
func generateLinuxUpdateScript(writer io.Writer, executablePaths []string, localBuildZipPath, lockName string, updateOperation config.UpdateOperation) error {
	template, err := template.New("linux-update-template").Parse(linuxReplaceBuildTemplate)
	if err != nil {
		return err
	}

	return template.Execute(writer, map[string]string{
		"ArchiveName":     filepath.Base(localBuildZipPath),
		"ExecutablePaths": csvify(executablePaths),
		"IsReplaceBuild":  getIsReplaceBuildTemplateValue(updateOperation),
		"LockName":        lockName,
	})
}

const linuxReplaceBuildTemplate = `
#!/bin/bash

set -e

ARCHIVE_NAME={{.ArchiveName}}
EXE_PATHS={{.ExecutablePaths}}
LOCKFILE="/tmp/{{.LockName}}.lock"
OLD_IFS="$IFS"

# Cleanup script at the end
function cleanup {
	flock -u 200
	exec 200>&-
	IFS="$OLD_IFS"
	rm -f $ARCHIVE_NAME
	rm -- "$0"
}
trap cleanup EXIT

echo "attempting to acquire update lock"
exec 200>$LOCKFILE
flock -n 200 || { echo "failed to acquire update lock another process is holding it"; exit 1; }
echo "update lock acquired"

{{if .IsReplaceBuild}}

IFS=","
for EXE_PATH in $EXE_PATHS
do
	echo "deleting existing executable: $EXE_PATH";
	sudo rm -f $EXE_PATH;
done

echo "unzipping the archive: /tmp/$ARCHIVE_NAME";
sudo unzip -o /tmp/$ARCHIVE_NAME -d /local/game && rm /tmp/$ARCHIVE_NAME;

echo "changing server permissions";
sudo chown -R gl-user-server:gl-user /local/game/*;

{{end}}

for EXE_PATH in $EXE_PATHS
do
	sudo chmod -R 774 $EXE_PATH;

	echo "killing running processes: $EXE_PATH";
	KILLED=$(sudo pkill -c -f "sudo -H -E -u gl-user-server $EXE_PATH");
	if [ "$KILLED" -gt 0 ]; then
		echo "killed $KILLED gameserver processes";
	else
		exit 1;
	fi
done
`
