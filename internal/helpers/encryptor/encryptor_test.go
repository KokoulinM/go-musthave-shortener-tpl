package encryptor

import (
	"testing"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/helpers"
	"github.com/gofrs/uuid"
)

func BenchmarkEncryptor_Encode(b *testing.B) {
	random, _ := helpers.GenerateRandom(16)

	encryptor, _ := New(random)

	userID, err := uuid.NewV4()
	if err != nil {
		return
	}

	b.ResetTimer()

	b.Run("encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			encryptor.Encode(userID.Bytes())
		}
	})
}

func BenchmarkEncryptor_Decode(b *testing.B) {
	random, _ := helpers.GenerateRandom(16)

	encryptor, _ := New(random)

	userID, err := uuid.NewV4()
	if err != nil {
		return
	}

	b.ResetTimer()

	b.Run("decode", func(b *testing.B) {
		encryptor.Decode(userID.String())
	})
}
