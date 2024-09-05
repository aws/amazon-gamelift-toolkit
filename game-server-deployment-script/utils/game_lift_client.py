import boto3

from models import ParsedBuild, ParsedFleet


class GameLiftClient:
    def __init__(self):
        self.client = boto3.client('gamelift')

    def create_build(self, parsed_build: ParsedBuild):
        callargs = dict(Name=parsed_build.name,
                        StorageLocation=parsed_build.storage_location,
                        Version=parsed_build.version,
                        OperatingSystem=parsed_build.operating_system,
                        Tags=parsed_build.tags,
                        ServerSdkVersion=parsed_build.server_sdk_version)
        return self.client.create_build(**{k: v for k, v in callargs.items() if v is not None})

    def create_fleet(self, parsed_fleet: ParsedFleet):
        callargs = dict(Name=parsed_fleet.name,
                        BuildId=parsed_fleet.build_id,
                        EC2InstanceType=parsed_fleet.ec2_instance_type,
                        RuntimeConfiguration=parsed_fleet.runtime_configuration,
                        Description=parsed_fleet.description,
                        EC2InboundPermissions=parsed_fleet.ec2_inbound_permissions,
                        NewGameSessionProtectionPolicy=parsed_fleet.new_game_session_protection_policy,
                        ResourceCreationLimitPolicy=parsed_fleet.resource_creation_limit_policy,
                        MetricGroups=parsed_fleet.metric_groups,
                        PeerVpcAwsAccountId=parsed_fleet.peer_vpc_aws_account_id,
                        PeerVpcId=parsed_fleet.peer_vpc_id,
                        FleetType=parsed_fleet.fleet_type,
                        InstanceRoleArn=parsed_fleet.instance_role_arn,
                        CertificateConfiguration=parsed_fleet.certificate_configuration,
                        Locations=parsed_fleet.locations,
                        Tags=parsed_fleet.tags,
                        ComputeType=parsed_fleet.compute_type,
                        AnywhereConfiguration=parsed_fleet.anywhere_configuration)
        return self.client.create_fleet(**{k: v for k, v in callargs.items() if v is not None})

    def delete_fleet(self, fleet_id: str):
        return self.client.delete_fleet(FleetId=fleet_id)

    def describe_alias(self, alias_id: str):
        return self.client.describe_alias(AliasId=alias_id)

    def describe_build(self, build_id: str):
        return self.client.describe_build(BuildId=build_id)

    def describe_fleet_attributes(self, fleet_id: str):
        return self.client.describe_fleet_attributes(FleetIds=[fleet_id])

    def describe_fleet_location_attributes(self, fleet_id: str, locations):
        return self.client.describe_fleet_location_attributes(FleetId=fleet_id,
                                                              Locations=locations)

    def describe_game_sessions(self, fleet_id: str):
        return self.client.describe_game_sessions(FleetId=fleet_id)

    def update_alias(self, alias_id: str, fleet_id: str):
        return self.client.update_alias(AliasId=alias_id,
                                        RoutingStrategy={
                                            'Type': 'SIMPLE',
                                            'FleetId': fleet_id,
                                            'Message': 'Updated via game_server_deployment script.'
                                        })
