package server

import (
	"context"
	"net/url"
	"regexp"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dauth/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var portSuffixRegex = regexp.MustCompile(`:[0-9]{2,5}$`)
var EmptyMetadata = metadata.New(nil)

type AuthenticatedServerStream struct {
	grpc.ServerStream
	AuthenticatedContext context.Context
}

func (s AuthenticatedServerStream) Context() context.Context {
	return s.AuthenticatedContext
}

func validateAuth(ctx context.Context, path string, authenticator dauth.Authenticator) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = EmptyMetadata
	}

	ctx, err := authenticator.Authenticate(ctx, path, url.Values(md), middleware.RealIP(peerFromContext(ctx), md))
	if err != nil {
		return ctx, err
	}
	return ctx, err
}
