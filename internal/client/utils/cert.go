package utils

import (
	"crypto/tls"

	"google.golang.org/grpc/credentials"

	"github.com/unbeman/ya-prac-go-second-grade/internal/client/config"
)

// LoadTLSCredentials loads client's certificate and private key for tls creds
func LoadTLSCredentials(cfg config.TLSConfig) (credentials.TransportCredentials, error) {
	clientCert, err := tls.LoadX509KeyPair(cfg.ClientCertPath, cfg.ClientPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	tlsCfg := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		InsecureSkipVerify: true,
	}
	return credentials.NewTLS(tlsCfg), nil
}
