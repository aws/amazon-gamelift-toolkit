# override local go binary (eg. use a specific version)
go = go

build-all: format vet test lint build-linux-amd64 build-windows-amd64

format:
	$(go) fmt ./...

vet:
	$(go) vet ./...

generate:
	$(go) generate ./...

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(go) build -o ./builds/fastbuild_linux_amd64 cmd/main.go

build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(go) build -o ./builds/fastbuild_windows_amd64.exe cmd/main.go

test: generate
	$(go) test ./... -cover -coverprofile coverage.html

coverage: test
	$(go) tool cover -html=coverage.html

lint:
	golangci-lint run

deps:
	# Download go modules
	$(go) mod download

test-deps:	
	# Download golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1
	golangci-lint --version

	# Code coverage tool
	$(go) get golang.org/x/tools/cmd/cover

	# Mock generator
	$(go) install github.com/matryer/moq@latest