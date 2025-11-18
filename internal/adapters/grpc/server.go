package grpcadapter

import (
	"context"
	"errors"
	"strings"

	"github.com/voita-pro/http-grpc-api/internal/domain"
	"github.com/voita-pro/http-grpc-api/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

// NewGRPCServer creates a grpc.Server and registers reflection and the AccountService.
func NewGRPCServer(auth domain.Authenticator, uc domain.Usecase, opts ...grpc.ServerOption) *grpc.Server {
	s := grpc.NewServer(opts...)
	pb.RegisterServicePBServer(s, NewHandler(auth, uc))
	reflection.Register(s)
	return s
}

func bearerFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no metadata from context")
	}
	vals := md.Get("authorization")
	if len(vals) == 0 {
		return "", errors.New("no authorization header")
	}
	auth := strings.TrimSpace(vals[0])
	if auth == "" {
		return "", errors.New("empty authorization header")
	}
	// Expecting: "Bearer <token>"
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization format")
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("unsupported authorization scheme")
	}
	tok := strings.TrimSpace(parts[1])
	if tok == "" {
		return "", errors.New("empty bearer token")
	}
	return tok, nil
}
