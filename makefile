# Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOINSTALL=$(GOCMD) install
GOGET=$(GOCMD) get
DOCKERCMD=docker-compose
BINARY_NAME=bin/go-lb

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

install:
	$(GOINSTALL)

run-server1:
	$(GORUN) cmd/server/main.go --type node --addr :8081

run-server2:
	$(GORUN) cmd/server/main.go --type node --addr :8082

run-lb:
	$(GORUN) cmd/server/main.go --type lb --addr :8080 --nodes localhost:8081,localhost:8082