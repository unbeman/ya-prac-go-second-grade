package model

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

type UserInput struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
}

type OTPInput struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type OTPOutput struct {
	SecretKey string `json:"secret_key"`
	AuthURL   string `json:"auth_url"`
}
