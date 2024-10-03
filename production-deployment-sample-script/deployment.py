import argparse
from models import ParsedBuild, ParsedFleet
from utils import GameLiftClient, utilities


# Parse out command line arguments and print help if needed
def _parse(args=None):
    parser = argparse.ArgumentParser(description='Script for deploying a new Build/Fleet version and deleting older resources.')
    parser.add_argument("-r", "--region", required=True, help='AWS Region, ex. us-west-2')
    parser.add_argument("-f", "--fleet-id", required=True, help='Existing FleetId being replaced, ex. fleet-1234-5678-90')
    parser.add_argument("-a", "--alias-id", required=True, help='Existing AliasId to update with new FleetId, ex. alias-1234-5678-90')
    parser.add_argument("-bj", "--build-json", required=True, help='Json file modeling the new Build resource, ex. ./new_build.json')
    parser.add_argument("-fj", "--fleet-json", required=True, help='Json file modeling the new Fleet resource, ex. ./new_fleet.json')
    parsed_args = vars(parser.parse_args())
    start(parsed_args)

# The main method of the script for performing the game server deployment
def start(args):
    # Validate input json files before moving on to network calls
    print("Opening Build json file: %s" % args['build_json'])
    parsed_build = ParsedBuild(args['build_json'])
    print("Opening Fleet json file: %s" % args['fleet_json'])
    parsed_fleet = ParsedFleet(args['fleet_json'])

    # Validate input FleetId and AliasId exist
    game_lift_client = GameLiftClient()
    previous_fleet_id = args['fleet_id']
    alias_id = args['alias_id']
    utilities.validate_fleet_exists(game_lift_client, previous_fleet_id)
    utilities.validate_alias_exists(game_lift_client, alias_id)

    # Create the Build based on input file
    print("\nCreating new Build with Name: %s" % parsed_build.name)
    create_build_response = game_lift_client.create_build(parsed_build)
    new_build_id = create_build_response['Build']['BuildId']
    print("Build %s created." % new_build_id)

    # Loop to wait until the Build is READY
    utilities.wait_for_build_to_be_ready(game_lift_client, new_build_id)
    print("Build %s is READY!" % new_build_id)

    # Create the Fleet based on inputs - Make sure to apply the new BuildId when creating the Fleet
    parsed_fleet.build_id = new_build_id
    print("\nCreating new Fleet with Name: %s, BuildId: %s" % (parsed_fleet.name, parsed_fleet.build_id))
    create_fleet_response = game_lift_client.create_fleet(parsed_fleet)
    new_fleet_id = create_fleet_response['FleetAttributes']['FleetId']
    print("Fleet %s created." % new_fleet_id)

    # Loop until Fleet is ACTIVE
    utilities.wait_for_fleet_to_be_active(game_lift_client, new_fleet_id)
    print("All Fleet locations on %s are ACTIVE!" % new_fleet_id)

    # ===============================================================================================================================
    # NOTE: If your new Fleet requires more than 1 instance, this is where to insert calls to UpdateFleetCapacity / PutScalingPolicy.
    # ===============================================================================================================================

    # Update the Alias with the new FleetId
    print("\nUpdating Alias %s from Fleet %s to new Fleet %s" % (alias_id, previous_fleet_id, new_fleet_id))
    update_alias_response = game_lift_client.update_alias(alias_id, new_fleet_id)

    # Loop on existing Fleet and wait for all GameSessions to end.  This can take a long time.
    utilities.wait_for_game_sessions_to_terminate(game_lift_client, previous_fleet_id)
    print("Fleet %s has 0 GameSessions, all traffic has transitioned to new Fleet %s." % (previous_fleet_id, new_fleet_id))

    # Delete the previous Fleet
    print("\nDeleting previous Fleet %s." % previous_fleet_id)
    game_lift_client.delete_fleet(previous_fleet_id)
    print("Deployment complete to Fleet %s!" % new_fleet_id)


# When called from the command line, will call into _parse, then start
if __name__ == "__main__":
    _parse()
