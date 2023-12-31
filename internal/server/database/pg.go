package database

import (
	"context"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
)

type pg struct {
	conn *gorm.DB
}

// NewPG returns the initialized pg object that implements Database interface.
func NewPG(cfg config.PG) (*pg, error) {
	db := &pg{}
	if err := db.connect(cfg.DSN); err != nil {
		return nil, err
	}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

// CreateUser inserts new given model.User if user.Login not occupied.
func (db *pg) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	result := db.conn.WithContext(ctx).Create(&user)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return user, fmt.Errorf("%w: login (%v)", ErrUserAlreadyExists, user.Login)
	}
	if result.Error != nil {
		return user, fmt.Errorf("%w: %v", ErrDatabase, result.Error)
	}
	return user, nil
}

// GetUserByLogin returns model.User with given login, if exists.
func (db *pg) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	user := model.User{}
	result := db.conn.WithContext(ctx).First(&user, "login = ?", login)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, fmt.Errorf("%w: login (%v)", ErrUserNotFound, login)
	}
	if result.Error != nil {
		return user, fmt.Errorf("%w: %v", ErrDatabase, result.Error)
	}
	return user, nil
}

// GetUserByID returns model.User with given user ID, if exists.
func (db *pg) GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user := model.User{}
	result := db.conn.WithContext(ctx).First(&user, userID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, fmt.Errorf("%w: ID (%v)", ErrUserNotFound, userID)
	}
	if result.Error != nil {
		return user, fmt.Errorf("%w: %v", ErrDatabase, result.Error)
	}
	return user, nil
}

// UpdateUser updates existed user.
func (db *pg) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	result := db.conn.WithContext(ctx).Model(&user).Updates(user)
	if result.Error != nil {
		return user, fmt.Errorf("%w: %v", ErrDatabase, result.Error)
	}
	return user, nil
}

// GetUserSecrets returns all user's credentials.
func (db *pg) GetUserSecrets(ctx context.Context, user model.User) ([]model.Credential, error) {
	var creds []model.Credential
	result := db.conn.WithContext(ctx).Find(&creds, "user_id = ?", user.ID)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, result.Error)
	}
	return creds, nil
}

// SaveUserSecrets upsert all user's credentials.
func (db *pg) SaveUserSecrets(ctx context.Context, user model.User) error {
	result := db.conn.WithContext(ctx).Save(&user)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDatabase, result.Error)
	}
	return nil
}

// DeleteUserSecrets delete all credentials for given user.
func (db *pg) DeleteUserSecrets(ctx context.Context, user model.User) error {
	result := db.conn.WithContext(ctx).Delete(&model.Credential{}, "user_id = ?", user.ID)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDatabase, result.Error)
	}
	return nil
}

// connect initialize database session connection instance with dsn.
func (db *pg) connect(dsn string) error {
	conn, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			TranslateError: true,
			Logger:         logger.Default.LogMode(logger.Info)},
	)
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

// migrate prepares database.
func (db *pg) migrate() error {
	tx := db.conn.Exec(fmt.Sprintf(`
	DO $$ BEGIN
		CREATE TYPE credential_type AS ENUM ('%v', '%v', '%v');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`, model.Login, model.Bank, model.Note))
	if tx.Error != nil {
		return tx.Error
	}
	return db.conn.AutoMigrate(
		&model.User{},
		&model.Credential{},
	)
}
