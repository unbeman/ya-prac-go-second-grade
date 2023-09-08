package services

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
)

var (
	ErrInvalid                = errors.New("invalid")
	ErrInvalidUserCredentials = errors.New("invalid login or password")
	ErrInvalidOTPToken        = errors.New("invalid otp token or user id")
	ErrInvalidJWTToken        = errors.New("invalid access token")
	ErrShouldVerifyOTP        = errors.New("verify 2fa token expected")
	ErrNoUser                 = errors.New("no user in context")
)

func GenStatusError(err error) error {
	var statusErr error
	switch {
	case errors.Is(err, ErrInvalid):
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrInvalidUserCredentials):
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrInvalidOTPToken):
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrInvalidJWTToken):
		statusErr = status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, ErrShouldVerifyOTP):
		statusErr = status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, database.ErrUserNotFound):
		statusErr = status.Error(codes.NotFound, err.Error())
	case errors.Is(err, database.ErrUserAlreadyExists):
		statusErr = status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, ErrNoUser):
		statusErr = status.Error(codes.Internal, err.Error())
	default:
		statusErr = status.Error(codes.Internal, err.Error())
	}
	return statusErr
}
