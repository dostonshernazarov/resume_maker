-include .env
export

CURRENT_DIR=$(shell pwd)
APP=api_gateway
CMD_DIR=./cmd

.PHONY: run
run:
	go run ${CMD_DIR}/app/main.go

.PHONY: proto-gen
proto-gen:
	./scripts/gen-proto.sh

.PHONY: swag-gen
swag-gen:
	swag init -g api/router.go -o api/docs
