Overview
--------
The purpose of this script is to upload a new version of a game server and transition player traffic from an existing
Fleet resource to a new Fleet resource by utilizing an Alias.  The Alias can be used as a GameSessionQueue destination
or call CreateGameSession using the AliasId rather than FleetId.

Prerequisites
-------------
1. Install AWS SDK for Python: https://aws.amazon.com/sdk-for-python/
1. Install python packages for unit testing:
   1. `pip install pytest`
   2. `pip install mock`

Limitations
-----------
As game architecture matures, multiple resources may be required, regions may expand, etc.  This script should serve
as a starting point, but will likely need to be customized to fit your specific deployment needs.
- Transferring traffic from 1 Fleet to another requires similar capacity being scaled out to handle large amounts of
  new traffic, this is not handled in the script.
- Similarly, once traffic is transitioned to the new Fleet, applying scaling policies is a best practice to reduce
  cost, this is not handled in the script.
- The script relies on existing resources as inputs.  This could be further automated utilizing tags to list and find
  existing resources to replace rather than requiring inputs, this is not an option in the script today.
- The script only supports Builds uploaded to your own S3 bucket.

Required Arguments
------------------
| Name | Short Syntax | Description |
| ---- | ------------ | ----------- |
| --region | -r | The AWS Region the Fleet and Alias resources are located in. |
| --fleet-id | -f | The FleetId of the existing Fleet being replaced. |
| --alias-id | -a | The AliasId of the existing Alias being replaced. |
| --build-json | -bj | A json file detailing the Build resource to be created. |
| --fleet-json | -fj | A json file detailing the Fleet resource to be created. |

Example Usage
-------------
1. Find existing Fleet to replace: `aws gamelift describe-fleet-attributes --region <region>` to find the FleetId.
   1. Learn more about DescribeFleetAttributes [here](https://docs.aws.amazon.com/gamelift/latest/apireference/API_DescribeFleetAttributes.html)
1. Find existing Alias to transfer traffic from: `aws gamelift list-aliases --region <region>` to find the AliasId.
   2. Learn more about ListAliases [here](https://docs.aws.amazon.com/gamelift/latest/apireference/API_ListAliases.html)
1. Create json skeleton for Build resource: `aws gamelift create-build --generate-cli-skeleton > new_build.json`
1. Create json skeleton for Fleet resource: `aws gamelift create-fleet --generate-cli-skeleton > new_fleet.json`
1. Fill in Build and Fleet attributes in new json files.  Ideally this only needs to be done once and future updates
   will have minimal fields to change.
   1. Learn more about CreateBuild [here](https://docs.aws.amazon.com/gamelift/latest/apireference/API_CreateBuild.html)
   1. Learn more about CreateFleet [here](https://docs.aws.amazon.com/gamelift/latest/apireference/API_CreateFleet.html)
1. Execute the script:
```
python3 ./game_server_deployment.py -r us-west-2 \
--fleet-id fleet-12345678-1234-1234-1234-12345678 \
--alias-id alias-12345678-1234-1234-1234-12345678 \
--build_json new_build.json \
--fleet_json new_fleet.json
```

How It Works
------------
1. Validate that the input Fleet and Alias exist with `describe-fleet-attributes` / `describe-alias`.
1. Call `create-build` based on input Build json file to create the new Build.
1. Loop calling `describe-build` until the Build is in state READY.
1. Call `create-fleet` based on input Fleet json file to create the new Fleet.
1. Loop calling `describe-fleet-attributes` until Fleet state is ACTIVE.
1. Loop calling `describe-fleet-location-attributes` until all Fleet Locations are ACTIVE.
1. Call `update-alias` to send traffic to the new Fleet.
1. Loop calling `describe-game-sessions` on the old Fleet until there are none remaining.
1. Call `delete-fleet` on the old Fleet.
