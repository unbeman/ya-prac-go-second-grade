package model

import uuid "github.com/satori/go.uuid"

type Credential struct {
	Base

	UserID uuid.UUID
	User   User

	Type     CredentialType `gorm:"type:credential_type;not null"`
	MetaData string
	Secret   []byte //encrypted by client's master password
	LocalID  uuid.UUID
}
