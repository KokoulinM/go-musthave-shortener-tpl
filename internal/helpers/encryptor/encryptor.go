// Package encryptor required for encryption and decryption
package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/google/uuid"
)

type Encryptor struct {
	// aesblock - represents an implementation of block cipher
	aesblock cipher.Block
	// key - encryption key
	key []byte
}

func New(key []byte) (*Encryptor, error) {
	enc := Encryptor{
		key: key,
	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	enc.aesblock = aesblock
	return &enc, nil
}

// Encode method encrypts the first block in src into dst and returns the hexadecimal encoding
func (e *Encryptor) Encode(value []byte) string {
	encrypted := make([]byte, aes.BlockSize)
	e.aesblock.Encrypt(encrypted, value)

	return hex.EncodeToString(encrypted)
}

// Decode method returns the bytes represented by the hexadecimal string s
func (e *Encryptor) Decode(value string) (string, error) {
	encrypted, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}

	decrypted := make([]byte, aes.BlockSize)
	e.aesblock.Decrypt(decrypted, encrypted)

	result, err := uuid.FromBytes(decrypted)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
