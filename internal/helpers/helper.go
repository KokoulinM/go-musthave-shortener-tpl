// Package helpers contains auxiliary functions
package helpers

import (
	"math/rand"
)

// randomInt generates a random number based on min- and max-
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// RandomString generates a random string based on len- string length
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

// GenerateRandom - generates a random sequence of bytes based on size
func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
