from unittest import TestCase
from unittest.mock import Mock
from utils import utilities
import pytest


class TestUtilities(TestCase):

    @pytest.fixture(autouse=True)
    def __inject_fixtures(self, mocker):
        self.mocker = mocker

    def setUp(self):
        self.client = Mock()

    def test_fleet_doesnt_exist(self):
        self.client.describe_fleet_attributes.return_value = {}
        with self.assertRaises(Exception):
            utilities.validate_fleet_exists(self.client, 'TestFleetId')

    def test_fleet_exists(self):
        self.client.describe_fleet_attributes.return_value = {'FleetAttributes': [{'FleetId': 'Anything'}]}
        utilities.validate_fleet_exists(self.client, 'TestFleetId')
        self.assertEqual(self.client.describe_fleet_attributes.call_count, 1)

    def test_alias_doesnt_exist(self):
        self.client.describe_alias.return_value = {}
        with self.assertRaises(Exception):
            utilities.validate_alias_exists(self.client, 'TestAliasId')

    def test_alias_exists(self):
        self.client.describe_alias.return_value = {'AliasId': 'Anything'}
        utilities.validate_alias_exists(self.client, 'TestAliasId')
        self.assertEqual(self.client.describe_alias.call_count, 1)

    def test_wait_for_build_to_be_ready(self):
        # Mock sleep to be 1 second
        self.mocker.patch("utils.utilities.BUILD_SLEEP", 1)
        # Return status NEW, then READY to exit the loop
        self.client.describe_build.side_effect = [{'Build': {'Status': 'NEW'}}, {'Build': {'Status': 'READY'}}]
        utilities.wait_for_build_to_be_ready(self.client, 'TestBuildId')
        # Verify 2 calls
        self.assertEqual(self.client.describe_build.call_count, 2)

    def test_wait_for_fleet_to_be_active_no_locations(self):
        # Mock sleep to be 1 second
        self.mocker.patch("utils.utilities.FLEET_SLEEP", 1)
        # Return status ACTIVATING, then ACTIVE, with no locations
        self.client.describe_fleet_attributes.side_effect = [{'FleetAttributes': [{'Status': 'ACTIVATING'}]},
                                                             {'FleetAttributes': [{'Status': 'ACTIVE'}]}]
        self.client.describe_fleet_location_attributes.side_effect = [{'LocationAttributes': []}]
        utilities.wait_for_fleet_to_be_active(self.client, 'TestFleetId')
        # Verify 2 calls on the fleet, 0 calls on fleet locations
        self.assertEqual(self.client.describe_fleet_attributes.call_count, 2)
        self.assertEqual(self.client.describe_fleet_location_attributes.call_count, 1)

    def test_wait_for_fleet_to_be_active_errors(self):
        # Mock sleep to be 1 second
        self.mocker.patch("utils.utilities.FLEET_SLEEP", 1)
        # Return status ACTIVATING, then ERROR, to trigger exception
        self.client.describe_fleet_attributes.side_effect = [{'FleetAttributes': [{'Status': 'ACTIVATING'}]},
                                                             {'FleetAttributes': [{'Status': 'ERROR'}]}]
        with self.assertRaises(Exception):
            utilities.wait_for_fleet_to_be_active(self.client, 'TestFleetId')

    def test_wait_for_fleet_to_be_active_with_locations(self):
        # Mock sleep to be 1 second
        self.mocker.patch("utils.utilities.FLEET_SLEEP", 1)
        # Return status ACTIVATING, then ACTIVE, with 3 locations
        self.client.describe_fleet_attributes.side_effect = [
            {'FleetAttributes': [{'Status': 'ACTIVATING', 'Locations': ['location-1', 'location-2', 'location-3']}]},
            {'FleetAttributes': [{'Status': 'ACTIVE', 'Locations': ['location-1', 'location-2', 'location-3']}]}]
        # Return each location going ACTIVE 1 call at a time
        self.client.describe_fleet_location_attributes.side_effect = [
            {'LocationAttributes': [{'LocationState': {'Location': 'location-1', 'Status': 'NEW'}}, {'LocationState': {'Location': 'location-2', 'Status': 'NEW'}}, {'LocationState': {'Location': 'location-3', 'Status': 'NEW'}}]},
            {'LocationAttributes': [{'LocationState': {'Location': 'location-1', 'Status': 'ACTIVE'}}, {'LocationState': {'Location': 'location-2', 'Status': 'NEW'}}, {'LocationState': {'Location': 'location-3', 'Status': 'NEW'}}]},
            {'LocationAttributes': [{'LocationState': {'Location': 'location-2', 'Status': 'ACTIVE'}}, {'LocationState': {'Location': 'location-3', 'Status': 'NEW'}}]},
            {'LocationAttributes': [{'LocationState': {'Location': 'location-3', 'Status': 'ACTIVE'}}]}]
        utilities.wait_for_fleet_to_be_active(self.client, 'TestFleetId')
        # Verify 2 calls on the fleet, 4 calls on the fleet locations
        self.assertEqual(self.client.describe_fleet_attributes.call_count, 2)
        self.assertEqual(self.client.describe_fleet_location_attributes.call_count, 4)

    def test_wait_for_game_sessions_to_terminate(self):
        # Mock sleep to be 1 second
        self.mocker.patch("utils.utilities.GAME_SESSION_SLEEP", 1)
        # Return 2 game sessions, 1 game session, and 0 game sessions
        self.client.describe_game_sessions.side_effect = [
            {'GameSessions': [{'GameSessionId': 1}, {'GameSessionId': 2}]},
                              {'GameSessions': [{'GameSessionId': 1}]},
                              {'GameSessions': []}]
        utilities.wait_for_game_sessions_to_terminate(self.client, 'TestFleetId')
        # Verify 3 calls were made
        self.assertEqual(self.client.describe_game_sessions.call_count, 3)
