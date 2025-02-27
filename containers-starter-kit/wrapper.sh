#!/bin/bash

# TODO: Set your hosting port here
PORT=7777
# TODO: Set your game server binary here
SERVER_BINARY_PATH="ServerBuild/yourgameserverbinaryhere"

echo "Start Go Wrapper in the background and register our port"
./SdkGoWrapper/gameliftwrapper $PORT &
WRAPPER_PID=$!

# NOTE: For Unreal, you can run the game server binary directly but if you do run the yourgameserver.sh, make sure to make both executable!
echo "Making sure we are able to execute the server binary"
chmod +x ./$SERVER_BINARY_PATH

echo "Start game server"
./$SERVER_BINARY_PATH

echo "Game server terminated, signal wrapper so it can call ProcessEnding()"
kill -SIGINT $WRAPPER_PID
echo "Sleep for 0.3 seconds to allow the wrapper to finish"
sleep 0.3