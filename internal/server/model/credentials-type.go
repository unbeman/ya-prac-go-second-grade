package model

import (
	"database/sql/driver"
	"fmt"
)

// CredentialType is type of credential.
type CredentialType string

const (
	Login CredentialType = "login"
	Bank  CredentialType = "bank"
	Note  CredentialType = "note"
)

func (st *CredentialType) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("CredentialType.Scan: failed to convert value")
	}
	*st = CredentialType(v)
	return nil
}

func (st CredentialType) Value() (driver.Value, error) {
	return string(st), nil
}
