package service

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/utils"
)

type OTPService struct {
	client pb.OtpServiceClient
	auth   *AuthService
}

func NewOTPService(conn grpc.ClientConnInterface, auth *AuthService) *OTPService {
	return &OTPService{client: pb.NewOtpServiceClient(conn), auth: auth}
}

func (s *OTPService) Generate() (string, string, error) {
	ctx := context.TODO()

	input := pb.OTPGenRequest{}

	out, err := s.client.OTPGenerate(ctx, &input)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			log.Errorf("Generate: status code %d, msg: %s", e.Code(), e.Message())
		} else {
			log.Errorf("Generate: %s", err)
		}
		return "", "", err
	}

	return out.SecretKey, out.AuthUrl, nil
}

func (s *OTPService) Validate(token string) error {
	ctx := context.TODO()

	input := pb.OTPValidateRequest{Token: token}

	out, err := s.client.OTPValidate(ctx, &input)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			log.Errorf("Validate: status code %d, msg: %s", e.Code(), e.Message())
		} else {
			log.Errorf("Validate: %s", err)
		}
		return err
	}

	s.auth.SetAccessToken(out.GetAccessToken())

	return nil
}

func (s *OTPService) Verify(token string) error {
	ctx := context.TODO()

	input := pb.OTPVerifyRequest{Token: token}

	out, err := s.client.OTPVerify(ctx, &input)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			log.Errorf("Verify: status code %d, msg: %s", e.Code(), e.Message())
		} else {
			log.Errorf("Verify: %s", err)
		}
		return err
	}

	s.auth.SetAccessToken(out.GetAccessToken())

	return nil
}

func (s *OTPService) Disable(login, password string) error {
	ctx := context.TODO()

	inputKey, err := utils.GetMasterKey(password, login)
	if err != nil {
		return fmt.Errorf("%w:%s", ErrInternal, err)
	}

	//todo: send hashed input key for comparing on the server side instead
	if string(inputKey) != string(s.auth.GetMaterKey()) {
		return fmt.Errorf("invalid login or master password")
	}

	input := pb.OTPDisableRequest{}

	_, err = s.client.OTPDisable(ctx, &input)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			log.Errorf("Disable: status code %d, msg: %s", e.Code(), e.Message())
		} else {
			log.Errorf("Disable: %s", err)
		}
		return err
	}

	return nil
}
