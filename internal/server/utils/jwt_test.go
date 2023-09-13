package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
)

func CreateKeyFile(t *testing.T, keyPath string) {

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(t, err)

	keyFile, err := os.Create(keyPath)
	require.NoError(t, err)

	defer keyFile.Close()

	pemKey := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(keyFile, pemKey)
	require.NoError(t, err)
}

func TestNewJWTManager(t *testing.T) {
	dir := t.TempDir()

	keyFile := dir + "/key.pem"
	CreateKeyFile(t, keyFile)

	noKeyFile := dir + "/nokey.pem"
	f, err := os.Create(noKeyFile)
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     config.JWT
		wantErr bool
	}{
		{
			name:    "good",
			cfg:     config.BuildTestJWTConfig(keyFile, config.AccessTokenExpiresInDefault),
			wantErr: false,
		},
		{
			name:    "no file",
			cfg:     config.BuildTestJWTConfig("", config.AccessTokenExpiresInDefault),
			wantErr: true,
		},
		{
			name:    "no key",
			cfg:     config.BuildTestJWTConfig(noKeyFile, config.AccessTokenExpiresInDefault),
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			_, err := NewJWTManager(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJWTManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestJWTManager_Generate(t *testing.T) {
	dir := t.TempDir()

	keyFile := dir + "/key.pem"
	CreateKeyFile(t, keyFile)

	jwtManager, err := NewJWTManager(config.BuildTestJWTConfig(keyFile, config.AccessTokenExpiresInDefault))
	require.NoError(t, err)

	cpy := *jwtManager

	jwtWithInvalidKey := &cpy
	shortKey, err := rsa.GenerateKey(rand.Reader, 256)
	require.NoError(t, err)

	jwtWithInvalidKey.privateKey = shortKey

	tests := []struct {
		name       string
		jwtManager IJWT
		user       model.User
		shouldF2A  bool
		wantErr    bool
	}{
		{
			name: "good",
			user: model.User{
				Base:       model.Base{ID: uuid.NewV4()},
				OtpEnabled: new(bool),
			},
			jwtManager: jwtManager,
			wantErr:    false,
		},
		{
			name: "bad key",
			user: model.User{
				Base:       model.Base{ID: uuid.NewV4()},
				OtpEnabled: new(bool),
			},
			jwtManager: jwtWithInvalidKey,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err = tt.jwtManager.Generate(tt.user, tt.shouldF2A)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestJWTManager_Verify(t *testing.T) {
	dir := t.TempDir()

	keyFile := dir + "/key.pem"
	CreateKeyFile(t, keyFile)

	jwtMan, err := NewJWTManager(config.BuildTestJWTConfig(keyFile, config.AccessTokenExpiresInDefault))
	require.NoError(t, err)

	goodToken, err := jwtMan.Generate(model.User{Base: model.Base{ID: uuid.NewV4()}, OtpEnabled: new(bool)}, false)
	require.NoError(t, err)

	cpy := *jwtMan

	anotherJwt := &cpy
	shortKey, err := rsa.GenerateKey(rand.Reader, 256)
	require.NoError(t, err)

	anotherJwt.privateKey = shortKey

	tests := []struct {
		name        string
		jwtManager  IJWT
		accessToken string
		want        *UserClaims
		wantErr     bool
	}{
		{
			name:        "good token",
			jwtManager:  jwtMan,
			accessToken: goodToken,
			wantErr:     false,
		},
		{
			name:        "no token",
			jwtManager:  jwtMan,
			accessToken: "",
			wantErr:     true,
		},
		{
			name:        "another rsa key",
			jwtManager:  anotherJwt,
			accessToken: goodToken,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.jwtManager.Verify(tt.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
