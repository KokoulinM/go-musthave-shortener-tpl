package storage

import (
	"errors"
	"strconv"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers"
)

type MockStorage storage

var testCount = 0

var data = []string{"https://go.dev", "https://mail.google.com", "https://practicum.yandex.ru", "http://localhost"}

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
	testCount += 1

	sl = string(helpers.RandomString(10) + "_test_" + strconv.Itoa(testCount))

	if s.Data == nil {
		s.Data = make(map[string]string)
	}

	s.Data[sl] = url

	return
}

func (s *MockStorage) Flush(c configs.Config) error {
	return nil
}

func (s *MockStorage) Load(c configs.Config) error {
	return nil
}
