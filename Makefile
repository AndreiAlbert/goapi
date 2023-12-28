GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_DIR=bin
SOURCE_DIR=src
BINARY_NAME=tuit

all: build
build: 
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v ./$(SOURCE_DIR)
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_DIR)/$(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v ./$(SOURCE_DIR)
	./$(BINARY_DIR)/$(BINARY_NAME)
fmt:
	go fmt ./...
deps:
	$(GOGET) ./$(SOURCE_DIR)/...

db:
	sudo docker-compose up -d
