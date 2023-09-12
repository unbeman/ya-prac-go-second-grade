package utils

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
)

type IJWT interface {
	Generate(user model.User, shouldF2A bool) (string, error)
	Verify(accessToken string) (*UserClaims, error)
}

type UserClaims struct {
	jwt.StandardClaims
	UserID     uuid.UUID `json:"user_id"`
	OtpEnforce bool      `json:"otp_enforce"`
}

type JWTManager struct {
	privateKey    *rsa.PrivateKey
	tokenDuration time.Duration
}

func NewJWTManager(cfg config.JWT) (*JWTManager, error) {
	rawPrivateKey, err := os.ReadFile(cfg.AccessTokenPrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(rawPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not parse key: %w", err)
	}

	return &JWTManager{
		privateKey:    privateKey,
		tokenDuration: cfg.AccessTokenExpiresIn,
	}, nil
}

func (m *JWTManager) Generate(user model.User, shouldF2A bool) (string, error) {
	now := time.Now().UTC()

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(m.tokenDuration).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
		UserID:     user.ID,
		OtpEnforce: shouldF2A && *user.OtpEnabled,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(m.privateKey)
	if err != nil {
		return "", fmt.Errorf("could not create signed token: %w", err)
	}
	return token, nil
}

func (m *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return m.privateKey.Public(), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
