package storage

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/unbeman/ya-prac-go-second-grade/internal/client/model"
)

type sqLite struct {
	conn *gorm.DB
}

func NewSqLiteDB(dbPath string) (*sqLite, error) {
	db := &sqLite{}
	conn, err := gorm.Open(
		sqlite.Open(dbPath),
		&gorm.Config{
			TranslateError: true,
			Logger:         logger.Default.LogMode(logger.Info)},
	)
	if err != nil {
		return nil, err
	}
	db.conn = conn
	err = db.migrate()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *sqLite) migrate() error {
	return db.conn.AutoMigrate(
		&model.Credential{},
	)
}

func (db *sqLite) AddCredential(ctx context.Context, cred *model.Credential) error {
	result := db.conn.WithContext(ctx).Create(cred)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return nil
}

func (db *sqLite) SaveCredential(ctx context.Context, cred model.Credential) (model.Credential, error) {
	log.Info("saving cred: ", cred.ID)
	result := db.conn.WithContext(ctx).Save(&cred)
	if result.Error != nil {
		return cred, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return cred, nil
}

func (db *sqLite) DeleteCredential(ctx context.Context, cred model.Credential) error {
	result := db.conn.WithContext(ctx).Delete(&cred)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return nil
}

func (db *sqLite) GetCredential(ctx context.Context, cred model.Credential) (model.Credential, error) {
	result := db.conn.WithContext(ctx).First(&cred)
	if result.Error != nil {
		return cred, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return cred, nil
}

func (db *sqLite) GetAllCredentials(ctx context.Context) ([]*model.Credential, error) {
	var creds []*model.Credential
	result := db.conn.WithContext(ctx).Order("meta_data ASC").Find(&creds)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return creds, nil
}

func (db *sqLite) SearchCredentials(ctx context.Context, search string) ([]*model.Credential, error) {
	var creds []*model.Credential

	//todo: avoid sql injection
	result := db.conn.WithContext(ctx).Where(fmt.Sprintf("meta_data LIKE '%%%s%%'", search)).Find(&creds)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return creds, nil
}

func (db *sqLite) GetAnyCredential(ctx context.Context) (model.Credential, error) {
	cred := model.Credential{}
	result := db.conn.WithContext(ctx).First(&cred)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return cred, fmt.Errorf("%w: %v", ErrNotFound, result.Error)
	}
	if result.Error != nil {
		return cred, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return cred, nil
}

func (db *sqLite) DeleteAll(ctx context.Context) error {
	result := db.conn.WithContext(ctx).Exec("DELETE FROM credentials")
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return nil
}
