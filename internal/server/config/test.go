package config

import "time"

func BuildTestJWTConfig(keyFile string, expDur time.Duration) JWT {
	return JWT{
		AccessTokenPrivateKeyFile: keyFile,
		AccessTokenExpiresIn:      expDur,
	}
}
