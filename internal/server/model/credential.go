package model

import uuid "github.com/satori/go.uuid"

// Credential describes entity for db credentials table.
type Credential struct {
	Base

	UserID uuid.UUID
	User   User

	Type     CredentialType `gorm:"type:credential_type;not null" sql:"type:ENUM('login', 'bank', 'note')"`
	MetaData string
	Secret   []byte //encrypted by client's master password
}
