module gameliftwrapper

replace aws/amazon-gamelift-go-sdk => ./gamelift-server-sdk

go 1.25.0

require aws/amazon-gamelift-go-sdk v0.0.0-00010101000000-000000000000

require (
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	golang.org/x/net v0.55.0 // indirect
)
