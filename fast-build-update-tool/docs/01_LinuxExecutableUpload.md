# Linux

Follow the instructions in this document if you are using GameLift instances hosted on Linux.

# Get Instance Access

This section describes how to get access to a remote GameLift instance.

If you are using an older GameLift SDK version (< 5), you will need to follow the instructions in the collapsed section titled "GameLift SDK < 5 Instructions" below.

## Generate an SSH Key

**[IMPORTANT]** You only need to do this step the first time you connect to an instance.

**Windows Users**: I recommend following the instructions in this section in Git Bash, rather than Powershell to avoid issues with Windows line endings and permissions.

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

## Copy the Public Key to the Instance

We need to start an SSM session on the instance, we can do this by getting access credentials via the AWS CLI

```sh
# Replace <fleet-id> with the id of your fleet
# Replace <instance-id> with the instance id of the instance you are connecting to
$ aws gamelift get-compute-access --fleet-id <fleet-id> --compute-name <instance-id>
{
    "Credentials": {
        "AccessKeyId": "<access-key-id>",
        "SecretAccessKey": "<secret-access-key>",
        "SessionToken": "<session-token>"
    }
}
```

We'll now need to use these credentials to connect to the instance. The easiest way to do this is using environment variables.

*IMPORTANT* Environment variables will override your AWS credentials. Any other AWS commands you run will use these credentials instead of your default credentials. You may want to run these commands in a new terminal session!

```sh
# Replace <access-key-id> with the access key returned by get-compute-access
# Replace <secret-access-key> with the secret access key returned by get-compute-access
# Replace <session-token> with the session token returned by get-compute-access
export AWS_ACCESS_KEY_ID=<access-key-id>
export AWS_SECRET_ACCESS_KEY=<secret-access-key>
export AWS_SESSION_TOKEN=<session-token>
```

Now we can start the SSM session! Once this command successfully completes we will be running commands on the remote instance.

```sh
aws ssm start-session --target <instance-id>
```

We need to enable SSH sessions with the public key we generated earlier (the output of the `ssh-keygen` command, it should start with `ssh-rsa`).

```sh
# Create an authorized_keys file for the gl-user-remote user.
$ sudo touch /home/gl-user-remote/.ssh/authorized_keys

# Append our public key to the authorized_keys file
# Replace <public-ssh-key with the output of ssh-keygen from earlier. It will be of the format `ssh-rsa $key`, you need to copy all of it!
$ echo "<public-ssh-key>" | sudo tee /home/gl-user-remote/.ssh/authorized_keys
```

## Open the SSH Port on the Instance

You will need to open up port 22 to your IP Address in order to remotely connect to the instance.

```sh
# Replace <fleet-id> with the id of your fleet
# Replace <ip-range> with the range of IPs you will be SSH-ing from. For a single IP address this will be "$your-public-ip/32" (eg. 127.0.0.1/32)
$ aws gamelift update-fleet-port-settings --fleet-id  <fleet-id>  --inbound-permission-authorizations "FromPort=22,ToPort=22,IpRange=<ip-range>,Protocol=TCP"
```

___

<details>
<summary>
GameLift SDK < 5 Instructions
</summary>

## GameLift SDK < 5

Only follow the instructions in this section if you are using a GameLift SDK version less than 5.

Most users can skip this section!

### Get Access Credentials

We can use the AWS CLI to get a `.pem` file used to get remote access to the instance

```sh
# Replace <fleet-id> with the id of your fleet
# Replace <instance-id> with the instance id of the instance you are connecting to
$ aws gamelift get-instance-access --fleet-id <fleet-id> --instance-id <instance-id> --query "InstanceAccess.Credentials.Secret" --output text > MyPrivateKey.pem
```

Depending on your SSH client you may need to change the permissions of the key.pem file in order to use it:

```sh
$ chmod 400 MyPrivateKey.pem
```

### Open the SSH Port on the Instance

You will need to open up port 22 to your IP Address in order to remotely connect to the instance.

```sh
# Replace <fleet-id> with the id of your fleet
# Replace <ip-range> with the range of IPs you will be SSH-ing from. For a single IP address this will be "$your-public-ip/32" (eg. 127.0.0.1/32)
$ aws gamelift update-fleet-port-settings --fleet-id  <fleet-id>  --inbound-permission-authorizations "FromPort=22,ToPort=22,IpRange=<ip-range>,Protocol=TCP"
```

</details>

___

## Upload and Replace an Executable

### Copy the Executable to the Instance

We'll need to look up the IP Address of the instance using the `describe-instances` command again.

```sh
# Replace <fleet-id> with the id of your fleet
# Replace <instance-id> with the instance id of the instance you are connecting to
$ aws gamelift describe-instances --fleet-id <fleet-id> --instance-id <instance-id>
{
    "Instances": [
        {
            "IpAddress": "127.0.0.1",
            # Abbreviated here for the sake of space
        }
    ]
}
```

Now we should be able to copy the executable up to the instance using `scp`. The default user for GameLift is `gl-user-remote`. We're copying these files to the `/tmp/` directory on the instance.

```sh
# Replace <instance-ip> with the IP address of the instance you are connecting to
# Replace <executable-file> with the server executable. If your executable requires additional files to run we recommend zipping them with the executable, and uploading the zip file.
$ scp -i MyPrivateKey.pem <executable-file> gl-user-remote@<instance-ip>:/tmp/
```

### Replacing the Server Executable in Place

We now have our server executable, and any related files on the remote instance, there are a few more steps required before we can use our new executable.

First let's actually connect to the instance using ssh.

```sh
# Replace <instance-ip> with the IP address of the instance you are connecting to
$ ssh -i MyPrivateKey.pem gl-user-remote@<instance-ip>
```

Now we are connected to the remote instance and can run commands. Previously, we copied our executable up to the `/tmp/` directory.
Next we'll have to change the owner of the executable, and copy it into the correct place. The `/local/game/` directory is where GameLift stores your server executable by default.

```sh
# Replace <executable> with the executable we wish to use. This can be an executable file or an entire folder.
$ sudo chown -R gl-user-server:gl-user /tmp/<executable>
$ sudo chmod -R 774 /tmp/<executable>

# Remove any files we'll be replacing in /local/game
# It's likely that we will want to remove everything in this folder, but we can selectively delete specific files that were uploaded with the previous build.
# At a minimum we need to replace the previous build's executable file!
$ rm -r /local/game/<executable>

# Replace <executable> with the executable we wish to use. This can be an executable file or an entire folder.
# If <executable> is a folder instead of a single file, you can use /tmp/<executable>/* instead of /tmp/<executable>
$ sudo mv /tmp/<executable> /local/game/
```

We have everything copied into place, now we just need to kill any lingering processes. 

```sh
# Use this command to list any running server processes
# Replace <executable-file-name> with the name of the executable used to start your game
$ ps -efww | grep "<executable-file-name>"

# Use this command to actually kill any of the processes listed out!
# Replace <executable-file-name> with the name of the executable used to start your game
$ sudo pkill -f <executable-file-name>
```

After the kill command is run, you should run the list processes command to make sure the previous processes have been stopped, and new ones are spun up. This may take a few moments of time depending on how long your server processes take to shut down and start up. Once new server processes are started you should be able to connect to them, and play like normal!
