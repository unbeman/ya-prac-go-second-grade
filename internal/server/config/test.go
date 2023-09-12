package config

import "time"

func BuildTestJWTConfig(keyFile string, expDur time.Duration) JWT {
	return JWT{
		AccessTokenPrivateKeyFile: keyFile,
		AccessTokenExpiresIn:      expDur,
	}
}

func BuildTestTLSConfig(cert, key string) TLS {
	return TLS{
		CertPath: cert,
		KeyPath:  key,
	}
}
