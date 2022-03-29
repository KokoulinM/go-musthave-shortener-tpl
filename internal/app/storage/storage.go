package storage

import (
	"errors"
	"sync"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers"
)

type Repository interface {
	LinkBy(sl string) (string, error)
	Save(url string) (sl string)
}

type Storage struct {
	data map[string]string
}

func (s *Storage) LinkBy(sl string) (string, error) {
	mu := sync.Mutex{}
	mu.Lock()

	link, ok := s.data[sl]
	if !ok {
		return link, errors.New("url not found")
	}

	mu.Unlock()

	return link, nil
}

func (s *Storage) Save(url string) (sl string) {
	mu := sync.Mutex{}
	mu.Lock()

	sl = string(helpers.RandomString(10))

	s.data[sl] = url

	mu.Unlock()
	return
}

func New() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}
