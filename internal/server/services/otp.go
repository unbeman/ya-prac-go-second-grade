package services

import (
	"context"
	"fmt"

	"github.com/pquerna/otp/totp"
	log "github.com/sirupsen/logrus"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/utils"
)

// OTP service implements proto pb.OtpServiceServer interface.
type OTP struct {
	pb.UnimplementedOtpServiceServer
	db         database.Database
	jwtManager *utils.JWTManager
	project    string
}

// NewOTPService setups new OTP instance.
func NewOTPService(cfg config.OTP, db database.Database, jwtManager *utils.JWTManager) *OTP {
	return &OTP{
		db:         db,
		project:    cfg.Project,
		jwtManager: jwtManager,
	}
}

// OTPGenerate returns generated data to setting two-factor auth.
func (s OTP) OTPGenerate(ctx context.Context, input *pb.OTPGenRequest) (*pb.OTPGenResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, GenStatusError(err)
	}
	log.Info("OTPGenerate", user)
	otpInfo, err := s.generate(ctx, user)
	if err != nil {
		return nil, GenStatusError(err)
	}

	out := pb.OTPGenResponse{}
	out.SecretKey = otpInfo.SecretKey
	out.AuthUrl = otpInfo.AuthURL

	return &out, nil
}
func (s OTP) OTPVerify(ctx context.Context, input *pb.OTPVerifyRequest) (*pb.OTPVerifyResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, GenStatusError(err)
	}
	log.Info("OTPVerify", user.Login, user.OtpVerified, *user.OtpEnabled)
	_, err = s.verify(ctx, input.Token, user)
	if err != nil {
		return nil, GenStatusError(err)
	}

	token, err := s.jwtManager.Generate(user, false)
	if err != nil {
		return nil, GenStatusError(err)
	}

	out := pb.OTPVerifyResponse{AccessToken: token}
	return &out, nil
}

func (s OTP) OTPValidate(ctx context.Context, input *pb.OTPValidateRequest) (*pb.OTPValidateResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, GenStatusError(err)
	}

	err = s.validate(input.Token, user.OtpSecret)
	if err != nil {
		return nil, GenStatusError(err)
	}

	token, err := s.jwtManager.Generate(user, false)
	if err != nil {
		return nil, GenStatusError(err)
	}

	out := pb.OTPValidateResponse{AccessToken: token}
	return &out, nil
}
func (s OTP) OTPDisable(ctx context.Context, input *pb.OTPDisableRequest) (*pb.OTPDisableResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, GenStatusError(err)
	}

	//todo: check input

	user, err = s.disable(ctx, user)
	if err != nil {
		return nil, GenStatusError(err)
	}

	out := pb.OTPDisableResponse{}
	return &out, nil
}

// todo: move to controller layer

func (s OTP) generate(ctx context.Context, user model.User) (model.OTPOutput, error) {
	var output model.OTPOutput

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.project,
		AccountName: user.Login,
	})

	if err != nil {
		return output, fmt.Errorf("error while gen otp: %w", err)
	}

	user.OtpSecret = key.Secret()
	user.OtpAuthUrl = key.URL()

	user, err = s.db.UpdateUser(ctx, user)
	if err != nil {
		return output, fmt.Errorf("error while gen otp: %w", err)
	}

	output.SecretKey = key.Secret()
	output.AuthURL = key.URL()

	return output, nil
}

func (s OTP) verify(ctx context.Context, token string, user model.User) (model.User, error) {
	valid := totp.Validate(token, user.OtpSecret)
	if !valid {
		return user, ErrInvalidOTPToken
	}

	*user.OtpEnabled = true
	*user.OtpVerified = true
	user, err := s.db.UpdateUser(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s OTP) validate(token, secret string) error {
	valid := totp.Validate(token, secret)
	if !valid {
		return ErrInvalidOTPToken
	}

	return nil
}

func (s OTP) disable(ctx context.Context, user model.User) (model.User, error) {
	*user.OtpEnabled = false
	user, err := s.db.UpdateUser(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func getUserFromContext(ctx context.Context) (model.User, error) {
	userValue := ctx.Value("auth-user")

	user, ok := userValue.(model.User)
	if !ok {
		return user, ErrNoUser
	}
	return user, nil
}
