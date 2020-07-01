#!/bin/bash
set -e

# Make the cleaning
rm -f rpc-server/*.pb.go
rm -f rpc-server/*.pb.gw.go
rm -f rpc-server/server_stub

# Compile proto files
protoc -Iproto/ -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I/usr/local/include --go_out=plugins=grpc:rpc-server --grpc-gateway_out=logtostderr=true:rpc-server proto/*.proto