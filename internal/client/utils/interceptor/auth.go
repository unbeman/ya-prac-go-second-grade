package interceptor

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/unbeman/ya-prac-go-second-grade/internal/client/service"
)

var ErrAuthExpected = errors.New("should login")

type AuthInterceptor struct {
	authService *service.AuthService
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor(authServie *service.AuthService, methods map[string]bool) *AuthInterceptor {
	//todo: schedule refresh token via login
	return &AuthInterceptor{authService: authServie, authMethods: methods}
}

// Unary returns a client interceptor to authenticate unary RPC
func (i *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if i.authMethods[method] {
			newCtx, err := i.attachToken(ctx)
			if err != nil {
				return err
			}
			return invoker(newCtx, method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Stream returns a client interceptor to authenticate stream RPC
func (i *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		if i.authMethods[method] {
			if i.authMethods[method] {
				newCtx, err := i.attachToken(ctx)
				if err != nil {
					return nil, err
				}
				return streamer(newCtx, desc, cc, method, opts...)
			}
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (i *AuthInterceptor) attachToken(ctx context.Context) (context.Context, error) {
	token := i.authService.GetAccessToken()
	if token == "" {
		return nil, ErrAuthExpected
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", token), nil
}
