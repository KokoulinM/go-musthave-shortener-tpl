package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"math/rand"
)

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

var aesGcm cipher.AEAD
var aesBlock cipher.Block
var nonce []byte

func Encode(userID string) (string, error) {
	src := []byte(userID)

	key, err := generateRandom(2 * aes.BlockSize) // ключ шифрования

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	aesBlock, err = aes.NewCipher(key)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	aesGcm, err = cipher.NewGCM(aesBlock)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	nonce, err = generateRandom(aesGcm.NonceSize())

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	dst := aesGcm.Seal(nil, nonce, src, nil) // зашифровываем
	fmt.Printf("encrypted: %x\n", dst)

	sha := hex.EncodeToString(dst)

	return sha, nil
}

func Decode(shaUserID string, userID *string) error {
	dst, err := hex.DecodeString(shaUserID)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}

	src, err := aesGcm.Open(nil, nonce, dst, nil)

	if err != nil {
		return err
	}

	*userID = string(src)

	return nil
}
