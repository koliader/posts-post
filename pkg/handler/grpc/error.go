package grpc_service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func errorResponse(code codes.Code, msg string) error {
	return status.Errorf(code, msg)
}
func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
