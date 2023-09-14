package model

// User describes user entity for db users table
type User struct {
	Base

	Login          string `gorm:"uniqueIndex;not null"`
	MasterKey2Hash string // master key twice hash

	OtpEnabled  *bool `gorm:"type:bool;default:false;"`
	OtpVerified *bool `gorm:"type:bool;default:false;"`
	OtpSecret   string
	OtpAuthUrl  string

	Credentials []Credential
}

type OTPOutput struct {
	SecretKey string `json:"secret_key"`
	AuthURL   string `json:"auth_url"`
}
