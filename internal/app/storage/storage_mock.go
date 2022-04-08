package storage

import (
	"errors"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers"
)

type MockStorage storage

var data = map[string]string{
	"GMWJGSAPGA": "https://go.dev",
	"TLMODYLUMG": "https://mail.google.com",
	"UQDIWWMNPP": "https://practicum.yandex.ru",
}

func (s *MockStorage) GenerateMockData() {
	for _, v := range data {
		s.Save(v)
	}
}

func (s *MockStorage) LinkBy(sl string) (string, error) {
	link, ok := s.Data[sl]
	if !ok {
		return link, errors.New("url not found")
	}

	return link, nil
}

func (s *MockStorage) Save(url string) (sl string) {
	sl = string(helpers.RandomString(10))

	if s.Data == nil {
		s.Data = make(map[string]string)
	}

	s.Data[sl] = url

	return
}
