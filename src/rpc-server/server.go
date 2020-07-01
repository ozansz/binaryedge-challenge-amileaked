package main

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type LeakServiceServerHandler struct {
	db *MongoDBConn
}

func (s *LeakServiceServerHandler) ListLeaks(req *empty.Empty, srv LeakService_ListLeaksServer) error {
	return status.Errorf(codes.Unimplemented, "method ListLeaks not implemented")
}
func (s *LeakServiceServerHandler) GetLeaksByEmail(req *GetLeaksByEmailRequest, srv LeakService_GetLeaksByEmailServer) error {
	return status.Errorf(codes.Unimplemented, "method GetLeaksByEmail not implemented")
}
func (s *LeakServiceServerHandler) GetLeaksByDomain(req *GetLeaksByDomainRequest, srv LeakService_GetLeaksByDomainServer) error {
	return status.Errorf(codes.Unimplemented, "method GetLeaksByDomain not implemented")
}
