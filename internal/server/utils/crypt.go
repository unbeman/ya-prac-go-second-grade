package utils

import "golang.org/x/crypto/bcrypt"

func HashToStore(key []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(key, bcrypt.DefaultCost)
}

func ValidateKey(storedHash, key []byte) error {
	return bcrypt.CompareHashAndPassword(storedHash, key)
}
