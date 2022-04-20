package auth

import (
	"context"
	"errors"
)

const (
	// AuthorizationMD authorization metadata name
	AuthorizationMD = "authorization"

	// GrpcGWCookieMD is cookie metadata name of GRPC in GRPC Gateway Request
	GrpcGWCookieMD = "grpcgateway-cookie"
)

var (
	// ErrMetadataMissing reports that metadata is missing in the incoming context.
	ErrMetadataMissing = errors.New("auth: could not locate request metadata")
	// ErrAuthorizationMissing reports that authorization metadata is missing in the incoming context.
	ErrAuthorizationMissing = errors.New("auth: could not locate authorization metadata")
	//ErrInvalidToken reports that the token is invalid.
	ErrInvalidToken = errors.New("auth: invalid token")
	// ErrMultipleAuthFound reports that too many authorization entries were found.
	ErrMultipleAuthFound = errors.New("auth: too many authorization entries")
)

// Authenticator defines the interface to perform the actual authentication of the request.
// Implementations should be fetch the required data from the context.Context object.
// GRPC specific data like `metadata` and `peer` is avilable on the context.
// Should return a new child node `context.Context` or `codes.Unauthenticated`when auth is
// lacking or `codes.PermissionDenied` when lacking permission.
type Authenticator interface {
	Authenticate(ctx context.Context) (context.Context, error)
}

// AuthenticatorFunc defines a function to perform authentication of requests.
// It returns a new child node `context.Context` or `codes.Unauthenticated`when auth is
// lacking or `codes.PermissionDenied` when lacking permission.
type AuthenticatorFunc func(ctx context.Context) (context.Context, error)

func (af AuthenticatorFunc) Authenticate(ctx context.Context) (context.Context, error) {
	return af(ctx)
}
