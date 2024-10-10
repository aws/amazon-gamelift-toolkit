import json


class ParsedBuild:
    def __init__(self, json_file: str):
        self.json_file = json_file
        try:
            with open(json_file, 'r') as file:
                build_json = json.load(file)
            # Mandatory inputs indexed directly to throw KeyError
            self.name = build_json['Name']
            self.storage_location = build_json['StorageLocation']
            # Optional inputs read with json get
            self.version = build_json.get('Version')
            self.operating_system = build_json.get('OperatingSystem')
            self.tags = build_json.get('Tags')
            self.server_sdk_version = build_json.get('ServerSdkVersion')
        except KeyError as e:
            print("Mandatory Build input: %s" % e)
            raise e
        except ValueError as e:
            print("Exception parsing Build json: %s" % e)
            raise e
