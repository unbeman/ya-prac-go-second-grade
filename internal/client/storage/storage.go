package storage

import (
	"context"
	"errors"

	"github.com/unbeman/ya-prac-go-second-grade/internal/client/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/model"
)

var ErrDB = errors.New("database error")
var ErrNotFound = errors.New("not found")

//todo: split for services

type IStorage interface {
	AddCredential(ctx context.Context, cred *model.Credential) error
	SaveCredential(ctx context.Context, cred model.Credential) (model.Credential, error)
	DeleteCredential(ctx context.Context, cred model.Credential) error
	DeleteAll(ctx context.Context) error
	GetCredential(ctx context.Context, cred model.Credential) (model.Credential, error)
	GetAnyCredential(ctx context.Context) (model.Credential, error)
	GetAllCredentials(ctx context.Context) ([]*model.Credential, error)
	SearchCredentials(ctx context.Context, search string) ([]*model.Credential, error)
}

func GetStorage(cfg config.AppConfig) (IStorage, error) {
	return NewSqLiteDB(cfg.DBFile)
}
