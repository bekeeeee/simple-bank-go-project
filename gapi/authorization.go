package gapi

import (
	"context"
	"github/bekeeeee/simplebank/token"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader     = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md.Get(authorizationHeader)
	if len(values) == 0 {

		return nil, status.Errorf(codes.Unauthenticated, "authorization header is not provided")

	}

	authHeader := values[0]

	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")

	}
	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return nil, status.Errorf(codes.Unauthenticated, "unsupported authorization type")

	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid access token")

	}
	return payload, nil

}
