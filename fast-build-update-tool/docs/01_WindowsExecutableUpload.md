# Windows

Follow the instructions in this document if you are using GameLift instances hosted on Windows.

# Getting Instance Access

This section describes how to get access to a remote GameLift instance.

If you are using an older GameLift SDK version (< 5), you will need to follow the instructions in the collapsed section titled "GameLift SDK < 5 Instructions" below.

## Get SSM Credentials

We first need to get AWS credentials that will be used to connect to the remote instance.

```sh
# Replace <fleet-id> with the id of your fleet
# Replace <instance-id> with the id of the instance you are connecting to
$ aws gamelift get-compute-access --fleet-id <fleet-id> --compute-name <instance-id>
{
    "Credentials": {
        "AccessKeyId": "<access-key-id>",
        "SecretAccessKey": "<secret-access-key>",
        "SessionToken": "<session-token>"
    }
}
```

## Connect to the Instance

We'll now need to use these credentials to connect to the instance. The easiest way to do this is through the use of environment variables.

*IMPORTANT* Environment variables will override your AWS credentials. Any other AWS commands you run will use these credentials instead of your default credentials. You may want to run these commands in a new terminal session!

```ps1
# Replace <access-key-id> with the access key returned by get-compute-access
# Replace <secret-access-key> with the secret access key returned by get-compute-access
# Replace <session-token> with the session token returned by get-compute-access
$Env:AWS_ACCESS_KEY_ID="<access-key-id>"
$Env:AWS_SECRET_ACCESS_KEY="<secret-access-key>"
$Env:AWS_SESSION_TOKEN="<session-token>"
```

Now we can start the SSM session! Once this command successfully completes we will be running commands on the remote instance.

```sh
# Replace <instance-id> with the id of the instance you are connecting to`
aws ssm start-session --target <instance-id>
```
___

<details>
<summary>
GameLift SDK < 5 Instructions
</summary>

## GameLift SDK < 5

Only follow the instructions in this section if you are using a GameLift SDK version less than 5.

Most users can skip this section!

### Open Port for RDP

First we'll need to open port 3389 in the fleet to allow us to RDP into the instance.

```ps1
# Replace <fleet-id> with the id of your fleet
# Replace <ip-range> with the range of IPs you will be SSH-ing from. For a single IP address this will be "$your-public-ip/32" (eg. 127.0.0.1/32)
$ aws gamelift update-fleet-port-settings --fleet-id <fleet-id>  --inbound-permission-authorizations "FromPort=3389,ToPort=3389,IpRange=<ip-range>,Protocol=TCP"
```

### RDP into the Instance

Next we'll need to get RDP credentials from GameLift using the `get-instance-access` command. This will return a username, password, and IP Address we can use to connect to the remote instance.

```ps1
# Replace <fleet-id> with the id of your fleet
# Replace <instance-id> with the id of the instance you are connecting to
$ aws gamelift get-instance-access ---fleet-id <fleet-id> --instance-id  <instance-id>
{
    "InstanceAccess": {
        "IpAddress": "127.0.0.1",
        "OperatingSystem": "WINDOWS_2016",
        "Credentials": {
            "UserName": "gl-user-remote",
            "Secret": "$password"
        }
        # Abbreviated here for the sake of space
    }
}
```

Next we can run the program "Remote Desktop Connection" in Windows, and start a connection to the remote instance.

Enter the IP address, username, and password returned by the previous command in the initial window that pops up.

Once you have access to the remote machine, open Powershell and run the following command to get admin privileges. This will open a new Powershell window.

```ps1
$ Start-Process powershell -Verb runAs
```

</details>

___

# Upload and Replace an Executable

## [Local Machine] Generate an SSH Key

**[IMPORTANT]** 
1. You only need to do this step the first time you connect to an instance.
2. Windows Users: I recommend following the instructions in this section in Git Bash, rather than Powershell to avoid issues with Windows line endings and permissions.
3. This section must happen on your local machine, not the remote server instance!

We'll need to first generate an SSH key for some of the later steps in this process.

We can quickly do this using the AWS CLI on your host machine (**not on the remote instance**).

```sh
$ aws ec2 create-key-pair --key-name gamelift-tool --region us-east-1 --query KeyMaterial --output text > MyPrivateKey.pem
```

Depending on your SSH client you may need to change the permissions of the key.pem file in order to use it later on:

```sh
$ chmod 400 MyPrivateKey.pem
```

We also need to grab the public key, and save this in a convenient spot. We'll need access to the contents of this key in a few moments.

```sh
$ ssh-keygen -y -f MyPrivateKey.pem
ssh-rsa <contents of public key>
```

On the host instance we'll also need to run a command to open port 1026 to our IP address so we can SSH to the instance later on. GameLift does not allow us to open the default SSH port `22` so we will use `1026` instead.

```sh
# Replace <fleet-id> with the id of your fleet
# Replace <ip-range> with the range of IPs you will be SSH-ing from. For a single IP address this will be "$your-public-ip/32" (eg. 127.0.0.1/32)
$ aws gamelift update-fleet-port-settings --fleet-id  <fleet-id>  --inbound-permission-authorizations "FromPort=1026,ToPort=1026,IpRange=<ip-range>,Protocol=TCP"
```

## [Remote Instance] Enable SSH on the Remote Windows Instance

In order to be able to upload an executable, we'll need to install OpenSSH on the remote instance.

First download OpenSSH, and extract it to the proper location:

```ps1
PS> [Net.ServicePointManager]::SecurityProtocol = "tls12, tls11, tls"
PS> Invoke-WebRequest -Uri "https://github.com/PowerShell/Win32-OpenSSH/releases/download/v8.9.1.0p1-Beta/OpenSSH-Win64.zip" -OutFile "C:\OpenSSH-Win64.zip"
PS> Expand-Archive C:\OpenSSH-Win64.zip -DestinationPath "C:\Program Files"
PS> Remove-Item -Path "C:\OpenSSH-Win64.zip"
```

Once we have OpenSSH in the Program Files folder, we need to run the install script, and start the OpenSSH process.

```ps1
PS> Set-Location -Path "C:\Program Files\OpenSSH-Win64\"
PS> powershell.exe -ExecutionPolicy Bypass -File install-sshd.ps1
PS> net start sshd
```

Update that SSH configuration to use port 1026 instead of 22.

```ps1
# Add port 1026 to the SSH config
PS> $sshConfig = "C:\ProgramData\ssh\sshd_config"
PS> $portLine = "Port 1026"
PS> $oldContent = Get-Content -Path $sshConfig
PS> $newContent = $portLine, $oldContent
PS> $newContent | Set-Content -Path $sshConfig
```

We also need to open port 1026 in the Windows firewall:

```ps1
PS> New-NetFirewallRule -Name sshd -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 1026
```

We need to add the Open SSH executables to the path so that we can copy files up to the instance using SCP:

```ps1
PS> setx PATH "$env:path;C:\Program Files\OpenSSH-Win64" -m
```

Now we need to add the public SSH key we obtained earlier using `ssh-keygen` to `authorized_keys`, so that we can SSH into this instance as `gl-user-server`.

```ps1
PS> New-Item -Path "C:\Users\gl-user-server\.ssh\" -ItemType Directory
# Replace <contents of public key> with the output of the `ssh-keygen` command from earlier
PS> Add-Content -Path "C:\Users\gl-user-server\.ssh\authorized_keys" -Value "ssh-rsa <contents of public key>"
```

Finally we can restart the SSH background process. Make sure that the `start` command finishes successfully before continuing!

```ps1
PS> net stop sshd
PS> net start sshd
```

## Replacing the Server Executable in Place 

Now that we have SSH access to the remote instance we need to copy up our executable.

This command should be run on your local machine, and will copy the executable to the `gl-user-server` user's home directory (`C:\Users\gl-user-server`).

```ps1
# Replace <executable> with the server executable. If your executable requires additional files to run we recommend zipping them with the exe file, and uploading the zip file.
$ scp -i MyPrivateKey.pem -P 1026 <executable> gl-user-server@<instance-ip>:<executable>
```

Now we'll need to SSH into the instance so we can finish this process

```ps1
# Replace <instance-ip> with the IP address of the instance you are connecting to
$ ssh -i MyPrivateKey.pem -p 1026 gl-user-server@<instance-ip>
```

If you zipped up your executable before copying it up, you can unzip it now using the `Expand-Archive` command.

We'll need to rename the running executable on the instance as Windows will not let you delete an `exe` that is currently running (moving a file is not an issue though).

```ps1
# Switch to Powershell :)
$ powershell

# Switch to the Game directory, this is where Gamelift stores your builds.
PS> Set-Location -Path  "C:\Game"

# Replace <your-executable> with the name of your game's exe
PS> Move-Item -Path "C:\Game\<your-executable>.exe" -Destination "C:\Game\<your-executable>-old.exe"
```

Now we need to copy our new build in-place to the Game folder.

```ps1
PS> Move-Item -Path "C:\Users\gl-user-server\<executable>" -Destination "C:\Game"
```

Finally we need to stop the running process. GameLift will automatically restart the game server process with your new executable in place.

```ps1
# Replace <executable-file-name> with your game's executable file without the file extension (eg. if your executable is `mygame.exe` you will call this command with `mygame`)
PS> Stop-Process -Name <executable-file-name> -Force
```

If you see everything is running successfully we can clean up the `*-old` files we moved around earlier to save space on the instance.
