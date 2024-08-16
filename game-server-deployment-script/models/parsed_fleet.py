import json


class ParsedFleet:
    def __init__(self, json_file: str):
        self.json_file = json_file
        try:
            with open(json_file, 'r') as file:
                fleet_json = json.load(file)
            # Mandatory inputs indexed directly to throw KeyError
            self.name = fleet_json['Name']
            self.build_id = fleet_json['BuildId']
            self.ec2_instance_type = fleet_json['EC2InstanceType']
            self.runtime_configuration = fleet_json['RuntimeConfiguration']
            # Optional inputs read with json get
            self.description = fleet_json.get('Description')
            self.ec2_inbound_permissions = fleet_json.get('EC2InboundPermissions')
            self.new_game_session_protection_policy = fleet_json.get('NewGameSessionProtectionPolicy')
            self.resource_creation_limit_policy = fleet_json.get('ResourceCreationLimitPolicy')
            self.metric_groups = fleet_json.get('MetricGroups')
            self.peer_vpc_aws_account_id = fleet_json.get('PeerVpcAwsAccountId')
            self.peer_vpc_id = fleet_json.get('PeerVpcId')
            self.fleet_type = fleet_json.get('FleetType')
            self.instance_role_arn = fleet_json.get('InstanceRoleArn')
            self.certificate_configuration = fleet_json.get('CertificateConfiguration')
            self.locations = fleet_json.get('Locations')
            self.tags = fleet_json.get('Tags')
            self.compute_type = fleet_json.get('ComputeType')
            self.anywhere_configuration = fleet_json.get('AnywhereConfiguration')
        except KeyError as e:
            print("Mandatory Fleet input: %s" % e)
            raise e
        except ValueError as e:
            print("Exception parsing Fleet json: %s" % e)
            raise e
