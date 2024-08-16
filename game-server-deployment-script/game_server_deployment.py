import argparse
import time
from datetime import datetime
from models import ParsedBuild, ParsedFleet
from utils import GameLiftClient

BUILD_SLEEP = 15
FLEET_SLEEP = 30
GAME_SESSION_SLEEP = 60

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
    validate_fleet_exists(game_lift_client, previous_fleet_id)
    validate_alias_exists(game_lift_client, alias_id)

    # Create the Build based on input file
    print("\nCreating new Build with Name: %s" % parsed_build.name)
    create_build_response = game_lift_client.create_build(parsed_build)
    new_build_id = create_build_response['Build']['BuildId']
    print("Build %s created." % new_build_id)

    # Loop to wait until the Build is READY
    wait_for_build_to_be_ready(game_lift_client, new_build_id)
    print("Build %s is READY!" % new_build_id)

    # Create the Fleet based on inputs
    print("\nCreating new Fleet with Name: %s" % args['fleet_name'])
    create_fleet_response = game_lift_client.create_fleet(parsed_fleet)
    new_fleet_id = create_fleet_response['FleetAttributes']['FleetId']
    print("Fleet %s created." % new_fleet_id)

    # Loop until Fleet is ACTIVE
    wait_for_fleet_to_be_active(game_lift_client, new_fleet_id)
    print("All Fleet locations on %s are ACTIVE!" % new_fleet_id)

    # ====================================================================================================
    # NOTE: Should consider scale here, the new fleet might need more than 1 instance or scaling policies.
    # ====================================================================================================

    # Update the Alias with the new FleetId
    print("\nUpdating Alias %s from Fleet %s to new Fleet %s" % (alias_id, previous_fleet_id, new_fleet_id))
    update_alias_response = game_lift_client.update_alias(alias_id, new_fleet_id)

    # Loop on existing Fleet and wait for all GameSessions to end.  This can take a long time.
    wait_for_game_sessions_to_terminate(game_lift_client, previous_fleet_id)
    print("Fleet %s has 0 GameSessions, all traffic has transitioned to new Fleet %s." % (previous_fleet_id, new_fleet_id))

    # Delete the previous Fleet
    print("\nDeleting previous Fleet %s." % previous_fleet_id)
    game_lift_client.delete_fleet(previous_fleet_id)
    print("Deployment complete!")

# Method for validating whether a fleet exists or not
def validate_fleet_exists(game_lift_client: GameLiftClient, _fleet_id: str):
    print("Validating %s exists..." % _fleet_id)
    existing_fleet_response = game_lift_client.describe_fleet_attributes(_fleet_id)
    if not existing_fleet_response['FleetAttributes']:
        raise Exception("Fleet %s was not found, exiting." % _fleet_id)

# Method for validating whether an alias exists or not
def validate_alias_exists(game_lift_client: GameLiftClient, _alias_id: str):
    print("Validating %s exists..." % _alias_id)
    existing_alias_response = game_lift_client.describe_alias(_alias_id)
    if not existing_alias_response:
        raise Exception("Alias %s was not found, exiting." % _alias_id)

# Method to wait for a build to have status READY before continuing
def wait_for_build_to_be_ready(game_lift_client: GameLiftClient, _build_id: str):
    new_build_state = 'NEW'
    while new_build_state != 'READY':
        print('Sleeping and describing Build %s...' % _build_id)
        time.sleep(BUILD_SLEEP)
        describe_build_response = game_lift_client.describe_build(_build_id)
        new_build_state = describe_build_response['Build']['Status']

# Method to wait for all fleet locations to have status ACTIVE before continuing
def wait_for_fleet_to_be_active(game_lift_client: GameLiftClient, _fleet_id: str):
    new_fleet_state = 'NEW'
    print('Sleeping and describing Fleet %s...' % _fleet_id)
    while new_fleet_state != 'ACTIVE' and new_fleet_state != 'ERROR':
        time.sleep(FLEET_SLEEP)
        describe_fleet_response = game_lift_client.describe_fleet_attributes(_fleet_id)
        new_fleet_state = describe_fleet_response['FleetAttributes'][0]['Status']
        print("%s: Fleet %s is still pending in status %s..." % (datetime.now(), _fleet_id, new_fleet_state))
    # If the fleet ended in ERROR, throw an exception and exit
    if new_fleet_state == 'ERROR':
        raise Exception("Fleet %s went into ERROR, exiting." % _fleet_id)
    # Check for Fleet locations and ensure those also go ACTIVE before continuing
    locations_response = game_lift_client.describe_fleet_location_attributes(_fleet_id, None)
    locations_remaining = [item.get('LocationState').get('Location') for item in locations_response['LocationAttributes']]
    while len(locations_remaining) != 0:
        time.sleep(FLEET_SLEEP)
        locations_response = game_lift_client.describe_fleet_location_attributes(_fleet_id, locations_remaining)
        print("%s: Waiting on locations %s to go ACTIVE: %s" %
              (datetime.now(), locations_remaining, locations_response['LocationAttributes']))
        for location in locations_response['LocationAttributes']:
            if location['LocationState']['Status'] == 'ACTIVE':
                locations_remaining.remove(location['LocationState']['Location'])

# Method to wait for all game sessions to terminate on a fleet before continuing
def wait_for_game_sessions_to_terminate(game_lift_client: GameLiftClient, _fleet_id: str):
    game_session_count = 1
    print("\nPolling previous Fleet %s for GameSessions." % _fleet_id)
    while game_session_count > 0:
        time.sleep(GAME_SESSION_SLEEP)
        describe_game_sessions_response = game_lift_client.describe_game_sessions(_fleet_id)
        game_session_count = len(describe_game_sessions_response['GameSessions'])
        print("%s: Fleet %s has %s or more GameSessions still running..." % (datetime.now(), _fleet_id, game_session_count))

# When called from the command line, will call into _parse, then start
if __name__ == "__main__":
    _parse()
