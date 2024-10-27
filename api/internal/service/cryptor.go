package service

import (
	"api/internal/entity"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

type userCryptor struct {
	key []byte
}

// NewUserCryptor initializes an Cryptor with a hex-encoded key
// the key parameter must be 32 bytes
func NewUserCryptor(key string) (*userCryptor, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, errors.New("key must be 32 bytes")
	}
	return &userCryptor{key: keyBytes}, nil
}

// Decrypt decrypts the provided encrypted user
func (e *userCryptor) Decrypt(user *entity.User) error {
	ciphertext, err := hex.DecodeString(user.Email)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	user.Email = string(ciphertext)

	return nil
}
