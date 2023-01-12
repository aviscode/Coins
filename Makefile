
BIN_DIR     := $(CURDIR)/bin
BIN_NAME   ?= Coins

default: build
fmt:
	go fmt ./...

vet:
	go vet ./...

mod:
	go mod tidy
	go mod verify

test:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go test -count=1 ./... -cover

build:mod fmt vet
	go build -o client ./Client-service/cmd
	go build -o server ./Pricing-service/cmd

