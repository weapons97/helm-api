
SHELL := /bin/bash

all: grpc vendor build

grpc:
	protoc -I. --go_out=plugins=grpc:pkg/protos --proto_path=./protos helm-api.proto
	protoc -I. --grpc-gateway_out=logtostderr=true,allow_delete_body=true:pkg/protos --proto_path=./protos helm-api.proto
	protoc -I. --openapiv2_out=logtostderr=true,allow_delete_body=true:swagger --proto_path=./protos helm-api.proto
	protoc -I. --doc_out=./doc --doc_opt=markdown,doc.md --proto_path=./protos helm-api.proto
build:
	go build -o bin/helm-api -mod=vendor
vendor:
	go mod tidy && go mod vendor