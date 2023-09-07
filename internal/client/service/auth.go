package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/utils"
)

var (
	ErrInternal          = errors.New("internal error")
	ErrLoginAlreadyExist = errors.New("login already exists")
)

type AuthService struct {
	client      pb.AuthServiceClient
	masterKey   []byte
	accessToken string
}

func NewAuthService(conn grpc.ClientConnInterface) *AuthService {
	client := pb.NewAuthServiceClient(conn)
	return &AuthService{client: client}
}

func (s *AuthService) GetMaterKey() []byte {
	return s.masterKey
}

func (s *AuthService) setMasterKey(key []byte) {
	s.masterKey = key
}

func (s *AuthService) GetAccessToken() string {
	return s.accessToken
}

func (s *AuthService) SetAccessToken(token string) {
	s.accessToken = token
}

func (s *AuthService) Register(login, masterPassword string) error {
	masterKey, err := utils.GetMasterKey(masterPassword, login)
	if err != nil {
		err = fmt.Errorf("%w : %s", ErrInternal, err)
		log.Error(err)
		return err
	}

	masterKeyHash, err := utils.GetMasterKeyHash(masterKey, masterPassword)
	if err != nil {
		err = fmt.Errorf("%w : %s", ErrInternal, err)
		log.Error(err)
		return err
	}

	ctx := context.TODO()

	input := pb.RegisterRequest{
		Login:   login,
		KeyHash: base64.StdEncoding.EncodeToString(masterKeyHash),
	}

	out, err := s.client.Register(ctx, &input)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			log.Errorf("Register: status code %d, msg: %s", e.Code(), e.Message())
			return ErrLoginAlreadyExist
		} else {
			log.Errorf("Register: %s", err)
		}
		return err
	}

	s.setMasterKey(s.masterKey)
	s.SetAccessToken(out.GetAccessToken())

	return nil
}

func (s *AuthService) Login(login, masterPassword string) error {
	masterKey, err := utils.GetMasterKey(masterPassword, login)
	if err != nil {
		err = fmt.Errorf("%w : %s", ErrInternal, err)
		log.Error(err)
		return err
	}

	masterKeyHash, err := utils.GetMasterKeyHash(masterKey, masterPassword)
	if err != nil {
		err = fmt.Errorf("%w : %s", ErrInternal, err)
		log.Error(err)
		return err
	}

	ctx := context.TODO()

	input := pb.LoginRequest{
		Login:   login,
		KeyHash: base64.StdEncoding.EncodeToString(masterKeyHash),
	}

	out, err := s.client.Login(ctx, &input)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			log.Errorf("Login: status code %d, msg: %s", e.Code(), e.Message())
		} else {
			log.Errorf("Login: %s", err)
		}
		return err
	}

	s.setMasterKey(masterKey)
	s.SetAccessToken(out.GetAccessToken())

	return nil
}
