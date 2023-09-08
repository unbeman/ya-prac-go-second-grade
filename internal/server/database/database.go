package database

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
)

type Database interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (model.User, error)

	GetUserSecrets(ctx context.Context, user model.User) ([]model.Credential, error)
	SaveUserSecrets(ctx context.Context, user model.User) error
	DeleteUserSecrets(ctx context.Context, user model.User) error
}

func GetDatabase(cfg config.PG) (Database, error) {
	return NewPG(cfg)
}
