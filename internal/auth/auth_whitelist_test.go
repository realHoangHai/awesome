package auth_test

import (
	"context"
	"google.golang.org/grpc"
	"testing"
)

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}

func TestAuthWhiteList(t *testing.T) {

}
