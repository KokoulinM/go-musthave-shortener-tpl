package storages

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/models"
)

type Repository struct {
	urls     models.ShortURLs
	filePath string
	baseURL  string
	usersURL map[models.UserID][]models.ShortURL
	mtx      sync.Mutex
}

type row struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
	User     string `json:"user"`
}

func FileRepository(ctx context.Context, filePath string, baseURL string) *Repository {
	repo := Repository{
		urls:     models.ShortURLs{},
		filePath: filePath,
		baseURL:  baseURL,
		usersURL: map[models.UserID][]models.ShortURL{},
	}

	cns, err := newConsumer(filePath)
	if err != nil {
		log.Printf("Error with reading file: %v\n", err)
	}
	defer cns.Close()

	reader := bufio.NewScanner(cns.file)
	for {
		ok, err := repo.readRow(reader)

		if err != nil {
			log.Printf("Error while parsing file: %v\n", err)
		}

		if !ok {
			break
		}
	}

	return &repo
}

func NewFileRepository(ctx context.Context, filePath string, baseURL string) handlers.Repository {
	return handlers.Repository(FileRepository(ctx, filePath, baseURL))
}

type Producer interface {
	//WriteEvent(event *storages)
	Close() error
}

type Consumer interface {
	//ReadEvent() (*storages, error)
	Close() error
}

type producer struct {
	file    *os.File
	write   *bufio.Writer
	encoder *json.Encoder
}

type consumer struct {
	file    *os.File
	read    *bufio.Reader
	decoder *json.Decoder
}

func (repo *Repository) AddURL(ctx context.Context, longURL, shortURL string, userID models.UserID) error {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	repo.urls[shortURL] = longURL
	err := repo.writeRow(longURL, shortURL, repo.filePath, userID)
	if err != nil {
		return errors.New("unexpected error when writing row")
	}

	repo.usersURL[userID] = append(repo.usersURL[userID], shortURL)

	return nil
}

func (repo *Repository) GetURL(ctx context.Context, sl models.ShortURL) (models.ShortURL, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	sl, ok := repo.urls[sl]
	if !ok {
		return "", errors.New("url not found")
	}

	return sl, nil
}

func (repo *Repository) GetUserURLs(ctx context.Context, userID models.UserID) ([]handlers.ResponseGetURL, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	var result []handlers.ResponseGetURL

	shortLinks := repo.usersURL[userID]

	for _, v := range shortLinks {
		result = append(result, handlers.ResponseGetURL{
			ShortURL:   fmt.Sprintf("%s/%s", repo.baseURL, v),
			OriginlURL: repo.urls[v],
		})
	}

	return result, nil
}

func (repo *Repository) Ping(ctx context.Context) error {
	return errors.New("not supported with filebase repository")
}

func newProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		write:   bufio.NewWriter(file),
		encoder: json.NewEncoder(file),
	}, nil
}

func newConsumer(filePath string) (*consumer, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0777)

	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		read:    bufio.NewReader(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (p *producer) Close() error {
	return p.file.Close()
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func (repo *Repository) readRow(reader *bufio.Scanner) (bool, error) {
	if !reader.Scan() {
		return false, reader.Err()
	}
	data := reader.Bytes()

	row := &row{}

	err := json.Unmarshal(data, row)

	if err != nil {
		return false, err
	}
	repo.urls[row.ShortURL] = row.LongURL
	repo.usersURL[row.User] = append(repo.usersURL[row.User], row.ShortURL)

	return true, nil
}

func (repo *Repository) writeRow(longURL, shortURL, filePath, userID string) error {
	p, err := newProducer(filePath)
	if err != nil {
		return err
	}

	data, err := json.Marshal(&row{
		LongURL:  longURL,
		ShortURL: shortURL,
		User:     userID,
	})
	if err != nil {
		return err
	}

	if _, err := p.write.Write(data); err != nil {
		return err
	}

	if err := p.write.WriteByte('\n'); err != nil {
		return err
	}

	return p.write.Flush()
}
