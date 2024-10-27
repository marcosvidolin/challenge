package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type Cryptor struct {
	key []byte
}

// NewCryptor initializes an Cryptor with a hex-encoded key
// the key parameter must be 32 bytes
func NewCryptor(key string) (*Cryptor, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, errors.New("key must be 32 bytes")
	}
	return &Cryptor{key: keyBytes}, nil
}

// Encrypt encrypts the provided val
func (e *Cryptor) Encrypt(val string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(val))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(val))

	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt decrypts the provided encrypted val
func (e *Cryptor) Decrypt(val string) (string, error) {
	ciphertext, err := hex.DecodeString(val)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
