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
	$(GORUN) main.go --type node --addr :8081

run-server2:
	$(GORUN) main.go --type node --addr :8082

run-lb:
	$(GORUN) main.go --type lb --addr :8080 --urls http://localhost:8081,http://localhost:8082