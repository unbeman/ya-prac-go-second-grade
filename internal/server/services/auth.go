package services

import (
	"context"
	"encoding/base64"
	"errors"

	uuid "github.com/satori/go.uuid"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/utils"
)

// Auth implements proto pb.AuthServiceServer interface.
type Auth struct {
	pb.UnimplementedAuthServiceServer
	db         database.Database
	JwtManager utils.IJWT
}

// NewAuthService setups new Auth instance.
func NewAuthService(db database.Database, jwtManager utils.IJWT) *Auth {
	return &Auth{db: db, JwtManager: jwtManager}
}

// Register creates new user.
func (s Auth) Register(ctx context.Context, input *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := s.register(ctx, input.GetLogin(), input.GetKeyHash())
	if err != nil {
		return nil, GenStatusError(err)
	}

	token, err := s.JwtManager.Generate(user, false)
	if err != nil {
		return nil, GenStatusError(err)
	}

	out := &pb.RegisterResponse{AccessToken: token}
	return out, nil
}

// Login authenticate user.
func (s Auth) Login(ctx context.Context, input *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.login(ctx, input.GetLogin(), input.GetKeyHash())
	if err != nil {
		return nil, GenStatusError(err)
	}

	token, err := s.JwtManager.Generate(user, *user.OtpEnabled)
	if err != nil {
		return nil, GenStatusError(err)
	}

	out := &pb.LoginResponse{AccessToken: token, Enforce_2FA: *user.OtpEnabled}
	return out, nil
}

// todo: move to controller layer

func (s Auth) login(ctx context.Context, login, key string) (model.User, error) {
	user, err := s.db.GetUserByLogin(ctx, login)
	if errors.Is(err, database.ErrUserNotFound) {
		return user, ErrInvalidUserCredentials
	}
	if err != nil {
		return user, err
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return user, err
	}

	decodedStoredKeyHash, err := base64.StdEncoding.DecodeString(user.MasterKey2Hash)
	if err != nil {
		return user, err
	}

	err = utils.ValidateKey(decodedStoredKeyHash, decodedKey)
	if err != nil {
		return user, ErrInvalidUserCredentials
	}
	return user, nil
}

func (s Auth) register(ctx context.Context, login, key string) (model.User, error) {
	var user model.User

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return user, err
	}

	hash, err := utils.HashToStore(decodedKey)
	if err != nil {
		return user, err
	}

	user.Login = login
	user.MasterKey2Hash = base64.StdEncoding.EncodeToString(hash)

	user, err = s.db.CreateUser(ctx, user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s Auth) GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error) {
	var user model.User

	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return user, err
	}

	return user, nil
}
