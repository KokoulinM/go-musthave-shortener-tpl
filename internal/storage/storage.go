package storage

import "errors"

type Storage struct {
	data map[string]string
}

type Repository interface {
	LinkBy(sl string) (string, error)
}

func (st *Storage) LinkBy(sl string) (string, error) {
	link, ok := st.data[sl]
	if !ok {
		return link, errors.New("url not found")
	}

	return link, nil
}

func New() *Storage {
	return &Storage{
		data: map[string]string{
			"123": "https://practicum.yandex.ru/",
		},
	}
}
