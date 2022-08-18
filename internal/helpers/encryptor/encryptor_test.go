package encryptor

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/helpers"
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
		_, err := encryptor.Decode(userID.String())
		if err != nil {
			fmt.Println(err)
		}
	})
}
