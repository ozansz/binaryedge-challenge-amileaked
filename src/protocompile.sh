#!/bin/bash
set -e

# Make the cleaning
rm -f rpc-server/*.pb.go
rm -f rest-server/*.pb.gw.go
rm -f rpc-server/server_stub

# Compile proto files
protoc -Iproto/ -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I/usr/local/include --go_out=plugins=grpc:rpc-server proto/*.proto
protoc -Iproto/ -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I/usr/local/include --go_out=plugins=grpc:rest-server --grpc-gateway_out=logtostderr=true:rest-server proto/*.proto