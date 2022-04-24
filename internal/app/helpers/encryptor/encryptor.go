package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type Encryptor struct {
	aesblock cipher.Block
	key      []byte
}

func New(key []byte) (*Encryptor, error) {
	enc := Encryptor{
		key: key,
	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		fmt.Errorf("error: %v\n", err)
		return nil, err
	}

	enc.aesblock = aesblock
	return &enc, nil
}

func (e *Encryptor) Encode(value []byte) string {
	encrypted := make([]byte, aes.BlockSize)
	e.aesblock.Encrypt(encrypted, value)

	return hex.EncodeToString(encrypted)
}

func (e *Encryptor) Decode(value string) (string, error) {
	encrypted, err := hex.DecodeString(value)
	if err != nil {
		fmt.Errorf("error: %v\n", err)
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
