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
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/storage"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/utils"
)

var (
	ErrInternal           = errors.New("internal error")
	ErrVault              = errors.New("can't decrypt vault")
	ErrLoginAlreadyExist  = errors.New("login already exists")
	ErrEnforceValidateOTP = errors.New("2Fa enabled; please, validate code")
)

type AuthService struct {
	client      pb.AuthServiceClient
	masterKey   []byte
	accessToken string
	vault       storage.IStorage
}

func NewAuthService(conn grpc.ClientConnInterface, vault storage.IStorage) *AuthService {
	client := pb.NewAuthServiceClient(conn)
	return &AuthService{client: client, vault: vault}
}

func (s *AuthService) GetMasterKey() []byte {
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
		err = fmt.Errorf("%w : %s", ErrInternal, err.Error())
		log.Error(err)
		return err
	}

	masterKeyHash, err := utils.GetMasterKeyHash(masterKey, masterPassword)
	if err != nil {
		err = fmt.Errorf("%w : %s", ErrInternal, err.Error())
		log.Error(err)
		return err
	}

	ctx := context.TODO()

	input := pb.RegisterRequest{
		Login:   login,
		KeyHash: base64.StdEncoding.EncodeToString(masterKeyHash),
	}

	err = s.vault.DeleteAll(ctx)
	if err != nil {
		err = fmt.Errorf("%w : %s", ErrInternal, err.Error())
		log.Error(err)
		return err
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

	s.setMasterKey(masterKey)
	s.SetAccessToken(out.GetAccessToken())

	return nil
}

//todo: мб разбить аутентификацию на сервере и офлайн вход на клиенте

func (s *AuthService) Login(login, masterPassword string) error {
	masterKey, err := utils.GetMasterKey(masterPassword, login)
	if err != nil {
		err = fmt.Errorf("%w : %s", ErrInternal, err.Error())
		log.Error(err)
		return err
	}

	err = s.checkVault(masterKey)
	if err != nil {
		log.Error(err)
		return err
	}

	masterKeyHash, err := utils.GetMasterKeyHash(masterKey, masterPassword)
	if err != nil {
		err = fmt.Errorf("%w : %s", ErrInternal, err.Error())
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

	if out.GetEnforce_2FA() {
		return ErrEnforceValidateOTP
	}

	return nil
}

func (s *AuthService) checkVault(masterKey []byte) error {
	ctx := context.TODO()
	cred, err := s.vault.GetAnyCredential(ctx)

	// no secrets in vault, may login
	if errors.Is(err, storage.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%w : %v", ErrInternal, err)
	}

	_, err = utils.Decrypt(masterKey, cred.Encrypted)
	if err != nil {
		return fmt.Errorf("%w : %v", ErrVault, err)
	}
	return nil
}
