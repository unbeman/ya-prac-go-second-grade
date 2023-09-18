package config

import "time"

// BuildTestJWTConfig setup jwt config for tests.
func BuildTestJWTConfig(keyFile string, expDur time.Duration) JWT {
	return JWT{
		AccessTokenPrivateKeyFile: keyFile,
		AccessTokenExpiresIn:      expDur,
	}
}

// BuildTestTLSConfig setup tls config for tests.
func BuildTestTLSConfig(cert, key string) TLS {
	return TLS{
		CertPath: cert,
		KeyPath:  key,
	}
}

// BuildTestOTPConfig setup otp config for tests.
func BuildTestOTPConfig() OTP {
	return OTP{
		Project: "test",
	}
}
