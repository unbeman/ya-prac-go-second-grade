package transport

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/services"
)

type AuthInterceptor struct {
	authService   *services.Auth
	noAuthMethods map[string]bool
	f2aMethods    map[string]bool
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(authService *services.Auth, noAuthMethods map[string]bool, f2aMethods map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{
		authService:   authService,
		noAuthMethods: noAuthMethods,
		f2aMethods:    f2aMethods,
	}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC.
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if _, ok := i.noAuthMethods[info.FullMethod]; !ok {

			user, err := i.authorize(ctx, info.FullMethod)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, "auth-user", user)

		}
		return handler(ctx, req)
	}
}

// authorize checks user by provided in metadata access token,
// prohibits method if user have enabled 2fa, but he didn't verify token.
func (i *AuthInterceptor) authorize(ctx context.Context, method string) (model.User, error) {
	var user model.User

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return user, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return user, status.Errorf(codes.Unauthenticated, "access token is not provided")
	}

	accessToken := values[0]

	claims, err := i.authService.JwtManager.Verify(accessToken)
	if err != nil {
		return user, fmt.Errorf("%w : %s", services.ErrInvalidJWTToken, err)
	}

	_, isF2a := i.f2aMethods[method]

	if claims.OtpEnforce && !isF2a {
		return user, fmt.Errorf("%w : %s", services.ErrShouldVerifyOTP, err)
	}

	user, err = i.authService.GetUserByID(ctx, claims.UserID)
	if errors.Is(err, database.ErrUserNotFound) {
		return user, fmt.Errorf("%w : %s", services.ErrInvalidJWTToken, err)
	}
	if err != nil {
		return user, services.GenStatusError(err)
	}

	return user, nil
}

type serverStreamWrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func newServerStreamWrapper(ss grpc.ServerStream, newCtx context.Context) *serverStreamWrapper {
	return &serverStreamWrapper{
		ServerStream: ss,
		ctx:          newCtx,
	}
}

func (w *serverStreamWrapper) Context() context.Context {
	return w.ctx
}
