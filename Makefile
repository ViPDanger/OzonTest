APP_NAME=MockXMLDaily
BUILD_DIR=bin
SRC=./cmd/main.go

HOST=:8080
MONGO_URI=mongodb://127.0.0.1:27017
MONGO_BDNAME=xml_daily
MONGO_USER=admin
MONGO_PASSWORD=admin
GRPC_HOST=:2020

.PHONY: all build clean run

all: build

build:
	mkdir -p $(BUILD_DIR)
	export PATH="$PATH:$(go env GOPATH)/bin"
	swag init -g cmd/main.go -o ./docs
	protoc --go_out=. --go-grpc_out=. ./proto/valCurs.proto
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC)
	docker compose build
run:
	 docker compose up -d  > /dev/null
	exec $(BUILD_DIR)/$(APP_NAME) --host=$(HOST) --mongoURI=$(MONGO_URI) --mongodbName=$(MONGO_BDNAME) --mongoUser=$(MONGO_USER) --mongoPassword=$(MONGO_PASSWORD) --grpcHost=$(GRPC_HOST)
clean:	
	rm -rf $(BUILD_DIR)
	docker compose down -v
