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
		s.Save("userID", v)
	}
}

func (s *MockStorage) LinkBy(userID, sl string) (string, error) {
	shortLinks, ok := s.Data[userID]
	if !ok {
		return "", errors.New("url not found")
	}

	url, ok := shortLinks[sl]
	if !ok {
		return "", errors.New("url not found")
	}

	return url, nil
}

func (s *MockStorage) Save(userID UserID, url string) ShortLink {
	testCount += 1

	sl := string(helpers.RandomString(10) + "_test_" + strconv.Itoa(testCount))

	currentUrls := ShortLinks{}

	if urls, ok := s.Data[userID]; ok {
		currentUrls = urls
	}

	currentUrls[sl] = url

	s.Data[userID] = currentUrls

	return sl
}

func (s *MockStorage) Flush(c configs.Config) error {
	return nil
}

func (s *MockStorage) Load(c configs.Config) error {
	return nil
}

func (s *MockStorage) LinksByUser(userID UserID) (ShortLinks, error) {
	return ShortLinks{}, nil
}
