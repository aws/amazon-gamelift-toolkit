from unittest import TestCase
from models.parsed_build import ParsedBuild


class TestParsedBuild(TestCase):
    def test_non_existent_file(self):
        with self.assertRaises(FileNotFoundError):
            ParsedBuild('./test/test_data/non_existent_file.json')

    def test_invalid_json(self):
        with self.assertRaises(ValueError):
            ParsedBuild('./test/test_data/invalid_json.json')

    def test_missing_mandatory_input(self):
        with self.assertRaises(KeyError):
            ParsedBuild('./test/test_data/build-missing-mandatory.json')

    def test_minimal_inputs(self):
        parsed_build = ParsedBuild('./test/test_data/build-minimal.json')
        self.assertEqual(parsed_build.name, 'TestBuild')
        self.assertEqual(parsed_build.storage_location,
                         {'Bucket': 'test-bucket', 'Key': 'test-key', 'RoleArn': 'test-role-arn'})
        self.assertIsNone(parsed_build.version)
        self.assertIsNone(parsed_build.operating_system)
        self.assertIsNone(parsed_build.tags)
        self.assertIsNone(parsed_build.server_sdk_version)

    def test_all_inputs(self):
        parsed_build = ParsedBuild('./test/test_data/build-full.json')
        self.assertEqual(parsed_build.name, 'TestBuild')
        self.assertEqual(parsed_build.version, 'TestVersion')
        self.assertEqual(parsed_build.storage_location,
                         {'Bucket': 'test-bucket', 'Key': 'test-key',
                          'RoleArn': 'test-role-arn', 'ObjectVersion': 'test-object-version'})
        self.assertEqual(parsed_build.operating_system, 'TestOperatingSystem')
        self.assertEqual(parsed_build.tags, [{'Key': 'test-tag-key-1', 'Value': 'test-tag-value-1'},
                                            {'Key': 'test-tag-key-2', 'Value': 'test-tag-value-2'}])
        self.assertEqual(parsed_build.server_sdk_version, 'TestServerSdkVersion')
