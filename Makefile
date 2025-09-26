GOOS=linux
CGO_ENABLED=0
BUILD_DIR=build
BINARY_NAME=FeeServer
PACKAGE=github.com/deepissue/fee_server
PWD=$(shell pwd)

.PHONY: all build clean run tidy publish build-mac

all: build

tidy:
	cd src && go mod tidy

build: tidy
	@mkdir -p $(BUILD_DIR)
	cd src && \
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) \
	go build -a -installsuffix cgo -o $(PWD)/$(BUILD_DIR)/$(BINARY_NAME) $(PACKAGE)

build-mac: tidy
	@mkdir -p $(BUILD_DIR)
	cd src && \
	CGO_ENABLED=$(CGO_ENABLED) \
	go build -a -installsuffix cgo -o $(PWD)/$(BUILD_DIR)/$(BINARY_NAME) $(PACKAGE)

run: tidy
	cd src && \
	go run main.go start --application fee --profile dev --config ../config/server.hcl --log.level=debug --log.path ../logs


run-prod: tidy
	cd src && \
	go run main.go start --application fee --profile prod --config ../config/server-prod.hcl --log.level=debug --log.path ../logs

publish: build
	echo 1
	# ssh root@47.111.5.245 mkdir -p /data/aid.pub/chain-monitor/logs
	# scp build/FeeServer config.hcl Dockerfile docker-compose.yml root@47.111.5.245:/data/aid.pub/chain-monitor

test: build
	scp -r config root@192.168.31.242:/opt/fee_server
	scp -r build/FeeServer root@192.168.31.242:/opt/fee_server

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
