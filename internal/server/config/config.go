package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
)

var (
	AddressDefault  = "0.0.0.0:8080"
	DSNDefault      = "postgresql://postgres:1211@localhost:5432/pkeep"
	CertPathDefault = "cert/server.crt"
	KeyPathDefault  = "cert/server.key"

	AccessTokenPrivateKeyFileDefault = "cert/jwt_key.pem"
	AccessTokenExpiresInDefault      = 60 * time.Minute

	ProjectDefault = "passkeeper"
)

type PG struct {
	DSN string `env:"POSTGRES_DSN"`
}

type JWT struct {
	AccessTokenPrivateKeyFile string        `env:"ACCESS_TOKEN_PRIVATE_KEY_FILE"`
	AccessTokenExpiresIn      time.Duration `env:"ACCESS_TOKEN_EXPIRED_IN"`
}

type OTP struct {
	Project string `env:"PROJECT"`
}

type TLS struct {
	CertPath string `env:"CERT_PATH"`
	KeyPath  string `env:"KEY_PATH"`
}

type ServerConfig struct {
	Address  string `env:"PK_SERVER_ADDRESS"`
	Postgres PG
	JWT      JWT
	OTP      OTP
	TLS      TLS
}

func (cfg *ServerConfig) parseEnv() error {
	return env.Parse(cfg)
}

func GetServerConfig() (ServerConfig, error) {
	cfg := ServerConfig{
		Address:  AddressDefault,
		Postgres: PG{DSN: DSNDefault},
		JWT: JWT{
			AccessTokenPrivateKeyFile: AccessTokenPrivateKeyFileDefault,
			AccessTokenExpiresIn:      AccessTokenExpiresInDefault,
		},
		OTP: OTP{
			Project: ProjectDefault,
		},
		TLS: TLS{
			CertPath: CertPathDefault,
			KeyPath:  KeyPathDefault,
		},
	}

	err := cfg.parseEnv()
	if err != nil {
		if err != nil {
			return cfg, fmt.Errorf("could not load config from env: %w", err)
		}
	}
	return cfg, nil
}
