package model

import (
	//"github.com/google/uuid"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Credential struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`

	MetaData string

	Type CredentialType `sql:"type:ENUM('login', 'bank', 'note')"`

	Encrypted []byte
	Decrypted []byte `gorm:"-"` // doesn't exist in credential table

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (c *Credential) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.NewV4()
	}
	return nil
}
