#!/bin/bash

echo "Building the Unreal GameLift Server SDK binaries for Amazon Linux 2023..."
docker buildx build --platform=linux/amd64 --output=./ --target=server .
echo "Done, now zipping the content.."
zip -r AL2023GameliftUE5sdk.zip lib*
echo "Done! Select Actions -> Download File and type /home/cloudshell-user/amazon-gamelift-toolkit/building-gamelift-server-sdk-for-unreal-engine-and-amazon-linux/AL2023GameliftUE5sdk.zip to download the binaries."
