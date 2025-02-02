package server

import (
	"context"

	"github.com/streamingfast/dauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryAuthChecker(check dauth.Authenticator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		childCtx, err := validateAuth(ctx, info.FullMethod, check)
		if err != nil {
			return nil, obfuscateErrorMessage(err)
		}

		return handler(childCtx, req)
	}
}

func StreamAuthChecker(check dauth.Authenticator) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		childCtx, err := validateAuth(ss.Context(), info.FullMethod, check)
		if err != nil {
			return obfuscateErrorMessage(err)
		}

		return handler(srv, AuthenticatedServerStream{ServerStream: ss, AuthenticatedContext: childCtx})
	}
}

func obfuscateErrorMessage(err error) error {
	if st, ok := status.FromError(err); ok {
		msg := st.Message()
		switch st.Code() {
		case codes.Internal, codes.Unavailable, codes.Unknown:
			msg = "error with authentication service, please try again later"
		}
		return status.Error(st.Code(), msg)
	}
	return status.Errorf(codes.Unauthenticated, "authentication: %s", err.Error())
}
