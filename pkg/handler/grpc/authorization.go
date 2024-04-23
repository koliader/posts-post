package grpc_service

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authBearer          = "bearer"
)

func (s *Server) authorizeUser(ctx context.Context) (*string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing auth header")
	}
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 1 {
		return nil, fmt.Errorf("invalid auth header format")
	}

	email := fields[0]
	return &email, nil
}
