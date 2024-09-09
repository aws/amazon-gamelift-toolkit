import time
from datetime import datetime
from utils import GameLiftClient


BUILD_SLEEP = 15
FLEET_SLEEP = 30
GAME_SESSION_SLEEP = 60


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
    print("\nPolling previous Fleet %s for GameSessions..." % _fleet_id)
    while game_session_count > 0:
        time.sleep(GAME_SESSION_SLEEP)
        describe_game_sessions_response = game_lift_client.describe_game_sessions(_fleet_id)
        game_session_count = len(describe_game_sessions_response['GameSessions'])
        print("%s: Fleet %s has %s or more GameSessions still running..." % (datetime.now(), _fleet_id, game_session_count))
