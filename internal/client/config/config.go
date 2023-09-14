package config

import (
	"fmt"

	"github.com/caarlos0/env/v8"
)

var (
	AddressDefault        = "0.0.0.0:8080"
	DBPathDefault         = "local.db"
	CertPathDefault       = "cert/RootCA.crt"
	PrivateKeyPathDefault = "cert/RootCA.key"
)

type TLSConfig struct {
	ClientCertPath       string `env:"TLS_CERT_PATH"`
	ClientPrivateKeyPath string `env:"TLS_KEY_PATH"`
}

type AppConfig struct {
	Address string `env:"PK_SERVER_ADDRESS"`
	DBFile  string `env:"SQLITE_PATH"`
	Certs   TLSConfig
}

func (cfg *AppConfig) parseEnv() error {
	return env.Parse(cfg)
}

func GetClientConfig() (AppConfig, error) {
	cfg := AppConfig{
		Address: AddressDefault,
		DBFile:  DBPathDefault,
		Certs: TLSConfig{
			ClientCertPath:       CertPathDefault,
			ClientPrivateKeyPath: PrivateKeyPathDefault,
		},
	}
	err := cfg.parseEnv()
	if err != nil {
		return cfg, fmt.Errorf("could not load config from env: %w", err)
	}
	return cfg, nil
}
