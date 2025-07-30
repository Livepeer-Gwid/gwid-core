package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"gwid.io/gwid-core/internal/config"
)

type EncryptionService struct {
	config *config.Config
}

func NewEncryptionService(config *config.Config) *EncryptionService {
	return &EncryptionService{
		config: config,
	}
}

func (s *EncryptionService) validateKey(key string) ([]byte, bool, int) {
	byteKey := []byte(key)

	if len(byteKey) == 32 {
		return byteKey, true, len(byteKey)
	} else {
		return byteKey, false, len(byteKey)
	}
}

func (s *EncryptionService) EncryptData(data []byte) (string, error) {
	key, isEncryptionKeyValid, size := s.validateKey(s.config.EncryptionKey)

	if !isEncryptionKeyValid {
		return "", fmt.Errorf("invalid encryption key of size %v", size)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, data, nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (s *EncryptionService) DecryptData(encryptedData string) (string, error) {
	key, isEncryptionKeyValid, size := s.validateKey(s.config.EncryptionKey)

	if !isEncryptionKeyValid {
		return "", fmt.Errorf("invalid encryption key of size %v", size)
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	data, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
