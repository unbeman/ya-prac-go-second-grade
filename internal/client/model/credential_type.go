package model

import "database/sql/driver"

type CredentialType string

const (
	Login CredentialType = "login"
	Bank                 = "bank"
	Note                 = "note"
)

func (st *CredentialType) Scan(value interface{}) error {
	*st = CredentialType(value.(string))
	return nil
}

func (st CredentialType) Value() (driver.Value, error) {
	return string(st), nil
}
