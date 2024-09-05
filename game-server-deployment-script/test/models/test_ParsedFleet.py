from unittest import TestCase
from models.parsed_fleet import ParsedFleet


class TestParsedBuild(TestCase):
    def test_non_existent_file(self):
        with self.assertRaises(FileNotFoundError):
            ParsedFleet('./test/test_data/non_existent_file.json')

    def test_invalid_json(self):
        with self.assertRaises(ValueError):
            ParsedFleet('./test/test_data/invalid_json.json')

    def test_missing_mandatory_input(self):
        with self.assertRaises(KeyError):
            ParsedFleet('./test/test_data/fleet-missing-mandatory.json')

    def test_minimal_inputs(self):
        parsed_fleet = ParsedFleet('./test/test_data/fleet-minimal.json')
        self.assertEqual(parsed_fleet.name, 'TestName')
        self.assertEqual(parsed_fleet.ec2_instance_type, 'TestInstanceType')
        self.assertEqual(parsed_fleet.runtime_configuration, {'ServerProcesses': [
            {'LaunchPath': 'TestLaunchPath', 'Parameters': 'test-parameter-1', 'ConcurrentExecutions': 1}]})
        self.assertIsNone(parsed_fleet.build_id)
        self.assertIsNone(parsed_fleet.description)
        self.assertIsNone(parsed_fleet.ec2_inbound_permissions)
        self.assertIsNone(parsed_fleet.new_game_session_protection_policy)
        self.assertIsNone(parsed_fleet.resource_creation_limit_policy)
        self.assertIsNone(parsed_fleet.metric_groups)
        self.assertIsNone(parsed_fleet.peer_vpc_aws_account_id)
        self.assertIsNone(parsed_fleet.peer_vpc_id)
        self.assertIsNone(parsed_fleet.fleet_type)
        self.assertIsNone(parsed_fleet.instance_role_arn)
        self.assertIsNone(parsed_fleet.certificate_configuration)
        self.assertIsNone(parsed_fleet.locations)
        self.assertIsNone(parsed_fleet.tags)
        self.assertIsNone(parsed_fleet.compute_type)
        self.assertIsNone(parsed_fleet.anywhere_configuration)

    def test_all_inputs(self):
        parsed_fleet = ParsedFleet('./test/test_data/fleet-full.json')
        self.assertEqual(parsed_fleet.name, 'TestName')
        self.assertEqual(parsed_fleet.build_id, 'TestBuildId')
        self.assertEqual(parsed_fleet.ec2_instance_type, 'TestInstanceType')
        self.assertEqual(parsed_fleet.runtime_configuration, {'ServerProcesses': [
            {'LaunchPath': 'TestLaunchPath', 'Parameters': 'test-parameter-1', 'ConcurrentExecutions': 1},
            {'LaunchPath': 'TestLaunchPath', 'Parameters': 'test-parameter-2', 'ConcurrentExecutions': 2}
            ], 'MaxConcurrentGameSessionActivations': 10, 'GameSessionActivationTimeoutSeconds': 10})
        self.assertEqual(parsed_fleet.description, 'TestDescription')
        self.assertEqual(parsed_fleet.ec2_inbound_permissions, [
            {'FromPort': 100, 'ToPort': 101, 'IpRange': '1.1.1.1/32', 'Protocol': 'UDP'},
            {'FromPort': 200, 'ToPort': 201, 'IpRange': '2.2.2.2/32', 'Protocol': 'UDP'}])
        self.assertEqual(parsed_fleet.new_game_session_protection_policy, 'TestProtectionPolicy')
        self.assertEqual(parsed_fleet.resource_creation_limit_policy,
                         {'NewGameSessionsPerCreator': 10, 'PolicyPeriodInMinutes': 10})
        self.assertEqual(parsed_fleet.metric_groups, ['test-metric-group-1', 'test-metric-group-2'])
        self.assertEqual(parsed_fleet.peer_vpc_aws_account_id, 'TestPeerVpcAwsAccountId')
        self.assertEqual(parsed_fleet.peer_vpc_id, 'TestPeerVpcId')
        self.assertEqual(parsed_fleet.fleet_type, 'TestFleetType')
        self.assertEqual(parsed_fleet.instance_role_arn, 'TestInstanceRoleArn')
        self.assertEqual(parsed_fleet.certificate_configuration, {'CertificateType': 'TestCertificateType'})
        self.assertEqual(parsed_fleet.locations, [{'Location': 'test-location-1'}, {'Location': 'test-location-2'}])
        self.assertEqual(parsed_fleet.tags, [{'Key': 'test-tag-key-1', 'Value': 'test-tag-value-1'},
                                             {'Key': 'test-tag-key-2', 'Value': 'test-tag-value-2'}])
        self.assertEqual(parsed_fleet.compute_type, 'TestComputeType')
        self.assertEqual(parsed_fleet.anywhere_configuration, {'Cost': 'TestCost'})
