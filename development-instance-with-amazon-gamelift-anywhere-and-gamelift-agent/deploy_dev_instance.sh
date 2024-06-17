#!/bin/bash

# TODO: Replace this with a globally unique name!
BUCKET_NAME="my-unique-bucket-name"


########## 1. CHECK THAT TOOLS ARE INSTALLED AND S3 BUCKET IS NOT OWNED BY SOMEONE ELSE, AND WE HAVE BUILT THE SAMPLE GAME SERVER ################

# Double check that java is installed and exit if not
if which java > /dev/null 2>&1; then
    echo "Java is installed"
else
    echo "Java not installed yet"
    exit 1
fi

# Double that Maven is installed and exit if not
if which java > /dev/null 2>&1; then
    echo "maven is installed"
else
    echo "Maven not installed yet"
    exit 1
fi

# Check that the sample game server is built
if [ ! -f "../../AmazonGameLiftSampleServerBinary.zip" ]; then
    echo "Sample game server not built yet."
    exit 1
fi

# Set the current region for the AWS CLI is us-east-1
echo "Setting region to us-east-1 for the AWS CLI"
aws configure set region us-east-1

# Check that the S3 bucket is not already owned by someone else
bucketstatus=$(aws s3api head-bucket --bucket "${BUCKET_NAME}" 2>&1)
if echo "${bucketstatus}" | grep 'Forbidden';
then
  echo "Bucket is already owned by someone else, edit deploy_dev_instance.sh to set a unique name"
  exit 1
elif echo "${bucketstatus}" | grep 'Bad Request';
then
  echo "Bucket name specified is less than 3 or greater than 63 characters"
  exit 1
else
  echo "You already have this bucket in your account, continue..."
fi


########## 2. CREATE THE S3 BUCKET, BUILD THE AGENT AND UPLOAD AGENT AND GAME SERVER BINARY TO S3 ################

# Create the S3 bucket
aws s3 mb s3://$BUCKET_NAME

# Build the Amazon GameLift Agent if it doesn't exist yet
agent_file="amazon-gamelift-agent/target/GameLiftAgent-1.0.jar"
if [ ! -f "$agent_file" ]; then
    git clone https://github.com/aws/amazon-gamelift-agent.git
    cd amazon-gamelift-agent/
    mvn clean compile assembly:single
    cd ..
else
    echo "Agent already built"
fi

# Copy the GameLift agent to the bucket
agent_file="amazon-gamelift-agent/target/GameLiftAgent-1.0.jar"
if [ ! -f "$agent_file" ]; then
    echo "Error: $agent_file does not exist"
    exit 1
fi
aws s3 cp "$agent_file" "s3://$BUCKET_NAME"

# Copy over the sample game server build we built before
aws s3 cp ../../AmazonGameLiftSampleServerBinary.zip s3://$BUCKET_NAME

########## 3. CREATE THE GAMELIFT RESOURCES ################

# Create the Amazon GameLift Anywhere location
LOCATION_NAME="custom-mygame-dev-location"
aws gamelift create-location --location-name $LOCATION_NAME

# Create the Amazon GameLift Anywhere fleet if it doesn't exist yet
FLEET_NAME="MyGame-Test-Fleet"
FLEET_ID=$(aws gamelift describe-fleet-attributes --query "FleetAttributes[?Name=='$FLEET_NAME'].FleetId" --output text 2>/dev/null)
if [ -z "$FLEET_ID" ]; then
    echo "Creating fleet: $FLEET_NAME"
    FLEET_ID=$(aws gamelift create-fleet --name $FLEET_NAME --compute-type ANYWHERE \
             --locations "Location=$LOCATION_NAME" \
             --runtime-configuration "ServerProcesses=[{LaunchPath=/local/game/GameLiftSampleServer,ConcurrentExecutions=1,Parameters=-logFile /local/game/logs/myserver1935.log -port 1935}]" \
             --anywhere-configuration Cost=0.2 \
             --query 'FleetAttributes.FleetId' --output text)
else
    echo "Fleet $FLEET_NAME already exists."
fi


########## 4. CREATE THE EC2 INSTANCE AND RELATED RESOURCES ################

# Only create the IAM and EC2 resources if we don't have an AmazonGameLiftDevInstance already
INSTANCE_ID=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=AmazonGameLiftDevInstance" "Name=instance-state-name,Values=running" --query "Reservations[*].Instances[*].InstanceId" --output text)
if [ -z "$INSTANCE_ID" ]; then

    # Create the IAM Role and Instance Profile for the EC2 instance
    aws iam create-role --role-name DevelopmentGameServerInstanceRole \
        --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}' \
        --description "Role for EC2 instance to run Amazon GameLift Agent" \
        --query 'Role.Arn'
    gamelift_policy=$(aws iam create-policy \
        --policy-name GameLiftFullAccess \
        --policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["gamelift:*"],"Resource":"*"}]}' \
        --query 'Policy.Arn' \
        --output text)
    aws iam attach-role-policy --role-name DevelopmentGameServerInstanceRole --policy-arn arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore
    aws iam attach-role-policy --role-name DevelopmentGameServerInstanceRole --policy-arn "$gamelift_policy"
    aws iam attach-role-policy --role-name DevelopmentGameServerInstanceRole --policy-arn arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy
    aws iam attach-role-policy --role-name DevelopmentGameServerInstanceRole --policy-arn arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess
    aws iam create-instance-profile --instance-profile-name GameLiftDevInstanceProfile
    # Wait for the instance profile to be created
    sleep 10
    aws iam add-role-to-instance-profile --role-name DevelopmentGameServerInstanceRole --instance-profile-name GameLiftDevInstanceProfile

    # Create a Security Group for the EC2 instance in the Default VPC
    SECURITY_GROUP_ID=$(aws ec2 create-security-group \
        --group-name game-server-sg \
        --description "Security group for the game server" \
        --vpc-id $(aws ec2 describe-vpcs --filters "Name=isDefault,Values=true" --query 'Vpcs[0].VpcId' --output text) \
        --query 'GroupId' --output text)

    # Allow inbound access for port 1935 for TCP
    aws ec2 authorize-security-group-ingress \
        --group-id $SECURITY_GROUP_ID \
        --protocol tcp \
        --port 1935 \
        --cidr 0.0.0.0/0 \
        --query 'SecurityGroupRules[0].SecurityGroupRuleId'

    # Create the EC2 instance and wait for it to start
    INSTANCE_ID=$(aws ec2 run-instances \
        --image-id resolve:ssm:/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64 \
        --instance-type m6i.large \
        --iam-instance-profile Name="GameLiftDevInstanceProfile" \
        --associate-public-ip-address \
        --security-group-ids $SECURITY_GROUP_ID \
        --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=AmazonGameLiftDevInstance}]' \
        --query 'Instances[0].InstanceId' \
        --output text)
    echo "Instance created with ID: $INSTANCE_ID"
    aws ec2 wait instance-running --instance-ids $INSTANCE_ID
else
    echo "Instance with ID $INSTANCE_ID already exists."
fi


########## 5. DEPLOY THE AGENT AND GAME SERVER BINARY TO THE EC2 INSTANCE AND CONFIGURE IT WITH SSM ################

# Configure and run the SSM command to install and start our game server
sed -i -e "s/your-fleet-id/$FLEET_ID/g" dev-game-server-setup-and-deployment.json
sed -i -e "s/your-bucket-name/$BUCKET_NAME/g" dev-game-server-setup-and-deployment.json

# Wait 15 seconds before sending the SSM command to make sure the SSM agent on the instance is ready
sleep 15

echo "EC2 instance is ready, sending SSM command to install and start the game server..."

aws ssm send-command --document-name "AWS-RunShellScript" \
--targets "Key=InstanceIds,Values=$INSTANCE_ID" \
--cli-input-json file://dev-game-server-setup-and-deployment.json \
--query 'Command.CommandId'

echo "All done! You should be able to start a game session in the next minute or so."

