package model

import "database/sql/driver"

type CredentialType string

const (
	Login CredentialType = "login"
	Bank  CredentialType = "bank"
	Note  CredentialType = "note"
)

func (st *CredentialType) Scan(value interface{}) error {
	*st = CredentialType(value.(string))
	return nil
}

func (st CredentialType) Value() (driver.Value, error) {
	return string(st), nil
}
