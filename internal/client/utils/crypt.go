package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/scrypt"
	"io"
)

//todo: error desc

func GetMasterKey(masterPassword, login string) ([]byte, error) {
	key, err := scrypt.Key([]byte(masterPassword), []byte(login), 32768, 8, 2, 64)
	return key, err
}

func GetMasterKeyHash(masterKey []byte, masterPassword string) ([]byte, error) {
	key, err := scrypt.Key(masterKey, []byte(masterPassword), 32768, 8, 2, 64)
	return key, err
}

// GetExtendedKey returns 32-bit key based on master key for data encryption.
func GetExtendedKey(masterKey []byte) ([]byte, error) {
	kr := hkdf.Expand(sha256.New, masterKey, nil)
	extendedKey := make([]byte, 32)
	_, err := io.ReadFull(kr, extendedKey)
	if err != nil {
		return nil, err
	}
	return extendedKey, nil
}

func Encrypt(masterKey, data []byte) ([]byte, error) {
	key, err := GetExtendedKey(masterKey) //todo: get once and save in memory
	if err != nil {
		return nil, err
	}

	ciphertext, err := aesEncrypt(key, data)
	if err != nil {
		return nil, err
	}

	encoded := make([]byte, base64.RawStdEncoding.EncodedLen(len(ciphertext)))
	base64.RawStdEncoding.Encode(encoded, ciphertext)

	return encoded, nil
}

func Decrypt(masterKey, encodedData []byte) ([]byte, error) {
	ciphertext := make([]byte, base64.RawStdEncoding.DecodedLen(len(encodedData)))
	_, err := base64.RawStdEncoding.Decode(ciphertext, encodedData)
	if err != nil {
		return nil, err
	}

	key, err := GetExtendedKey(masterKey) //todo: get once and save in memory
	if err != nil {
		return nil, err
	}

	data, err := aesDecrypt(key, ciphertext)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func aesEncrypt(aesKey []byte, plaintext []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func aesDecrypt(aesKey []byte, ciphertext []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		log.Error("NewCipher ", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Error("NewGCM ", err)
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Error("gcm.Open ", err)
		return nil, err
	}

	return plaintext, nil
}
