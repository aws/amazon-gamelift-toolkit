
# How to Manually Update a Server Executable

This guide describes how to manually update a server executable in-place on existing fleet instances without going through the typical GameLift build and fleet creation process.

### **NOTE** This process is not a requirement of using this CLI tool!

## Prerequisites

1. You will need access to a server executable for your game (and any files it depends on).
1. To take advantage of this process you must have a pre-existing GameLift fleet that runs on Managed EC2 instances.
1. You will need to have the [AWS CLI](https://aws.amazon.com/cli/) installed on your local machine.
    * You must have valid IAM access credentials. There are a number of ways to set this up, [they are outlined here](https://docs.aws.amazon.com/cli/v1/userguide/cli-chap-authentication.html).
    * If you run into any issues with access or permissions during this process, it is recommended that your IAM role or user is allowed to take the actions listed below. You can find information on how to change permissions for a user [here](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_change-permissions.html).
        * `gamelift:ListFleets`
        * `gamelift:DescribeFleetAttributes`
        * `gamelift:DescribeInstances`
        * `gamelift:GetComputeAccess`
        * `gamelift:UpdateFleetPortSettings`
        * `ec2:CreateKeyPair`
1. For users with Windows on their local machine:
    * You will need an SSH client installed. Git for Windows comes bundled with OpenSSH. You can also install OpenSSH separately.
    * Most commands outlined here are platform agnostic, but some of them require bash to work properly. Any such commands have been tested and function properly in Git Bash which comes bundled with Git for Windows. You may also be able to run equivalent commands in Powershell.
1. GameLift uses [SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-agent.html) to manage remote instance connections. **[In order to use SSM you will need to install the SSM CLI plugin from Amazon. You can find instructions to do this here.](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)**

## Finding an Instance to Update

### Get the Fleet ID

The first thing we will need to find is the unique id of the fleet you would like to update.

NOTE: If you do not see your fleet in the results you may not be querying the proper AWS region (you can override this using the `--region` parameter).

We'll want to use one of the `FleedId`s from these results.

```sh
# List available fleets
$ aws gamelift list-fleets
{
    "FleetIds": [
        "fleet-a1b2c3d4-5678-90ab-cdef-EXAMPLE11111",
        "fleet-a1b2c3d4-5678-90ab-cdef-EXAMPLE33333",
        "fleet-a1b2c3d4-5678-90ab-cdef-EXAMPLE44444"
    ]
}
```

```sh
# Replace <fleet-id> with one or more of the fleet ids returned by the `list-fleets` command separated by a comma
$ aws gamelift describe-fleet-attributes --fleet-ids <fleet-id>
{
    "FleetAttributes": [
        {
            "FleetId": "fleet-a1b2c3d4-5678-90ab-cdef-EXAMPLE11111",
            "Name": "my fleet",
            "BuildId": "build-11111111-1111-1111-1111-111111111111"
            # Abbreviated here for the sake of space
        }
    ]
}
```

### Get an Instance ID and Operating System

Now that we have our fleet ID we need to find an instance to actually update. Make sure that you take note of both the `InstanceId`, and `OperatingSystem` here, this will determine which guide you will need to follow later.

```sh
# Replace <fleet-id> with the id of your fleet
$ aws gamelift describe-instances --fleet-id <fleet-id>
{
    "Instances": [
        {
            "FleetId": "fleet-a1b2c3d4-5678-90ab-cdef-EXAMPLE11111",
            "InstanceId": "i-01111111111111111",
            "OperatingSystem": "AMAZON_LINUX_2",
            # Abbreviated here for the sake of space
        }
    ]
}
```

## Getting Access to an Instance and Updating the Executable

We now have all of the information we need to proceed:

* `FleetId`
* `InstanceId`
* `OperatingSystem`

If you do not know any of these values please revisit the previous steps!

At this point you must choose the correct guide to follow based on which `OperatingSystem` the GameLift fleet is using.

Further instructions are found here:

* [Windows Server](./01_WindowsExecutableUpload.md)
* [Linux Server](./01_LinuxExecutableUpload.md)


