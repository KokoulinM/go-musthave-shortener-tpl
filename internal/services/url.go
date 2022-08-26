package services

import (
	"context"
	"fmt"
	"net"

	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/handlers"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/models"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/shortener"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/workers"
)

type RepositoryInterface interface {
	AddURL(ctx context.Context, longURL models.LongURL, shortURL models.ShortURL, user models.UserID) error
	AddURLs(ctx context.Context, user models.UserID, urls ...handlers.RequestGetURLs) ([]handlers.ResponseGetURLs, error)
	DeleteURLs(ctx context.Context, user models.UserID, urls ...string) error
	GetURL(ctx context.Context, shortURL models.ShortURL) (models.ShortURL, error)
	GetUserURLs(ctx context.Context, user models.UserID) ([]handlers.ResponseGetURL, error)
	GetStates(ctx context.Context) (handlers.ResponseStates, error)
	Ping(ctx context.Context) error
}

type URLService struct {
	repo    RepositoryInterface
	baseURL string
	wp      *workers.WorkerPool
	subnet  *net.IPNet
}

func New(repo RepositoryInterface, baseURL string, wp *workers.WorkerPool, subnet *net.IPNet) *URLService {
	return &URLService{
		repo:    repo,
		baseURL: baseURL,
		wp:      wp,
		subnet:  subnet,
	}
}

func (us *URLService) GetURL(ctx context.Context, userID models.UserID) (string, error) {
	return us.repo.GetURL(ctx, userID)
}

func (us *URLService) CreateURL(ctx context.Context, longURL models.LongURL, user models.UserID) (string, error) {
	shortURL := shortener.ShorterURL(longURL)
	err := us.repo.AddURL(ctx, longURL, shortURL, user)
	return fmt.Sprintf("%s/%s", us.baseURL, shortURL), err
}

func (us *URLService) GetUserURLs(ctx context.Context, userID models.UserID) ([]handlers.ResponseGetURL, error) {
	return us.repo.GetUserURLs(ctx, userID)
}

func (us *URLService) Ping(ctx context.Context) error {
	return us.repo.Ping(ctx)
}

func (us *URLService) CreateBatch(ctx context.Context, urls []handlers.RequestGetURLs, userID models.UserID) ([]handlers.ResponseGetURLs, error) {
	return us.repo.AddURLs(ctx, userID, urls...)
}

func (us *URLService) DeleteBatch(urls []string, userID models.UserID) {
	var sliceData [][]string
	for i := 10; i <= len(urls); i += 10 {
		sliceData = append(sliceData, urls[i-10:i])
	}
	rem := len(urls) % 10
	if rem > 0 {
		sliceData = append(sliceData, urls[len(urls)-rem:])
	}
	for _, item := range sliceData {
		func(taskData []string) {
			us.wp.Push(func(ctx context.Context) error {
				err := us.repo.DeleteURLs(ctx, userID, taskData...)
				return err
			})
		}(item)
	}
}

func (us *URLService) GetStates(ctx context.Context, ip net.IP) (bool, handlers.ResponseStates, error) {
	if us.subnet == nil || !us.subnet.Contains(ip) {
		return false, handlers.ResponseStates{}, nil
	}
	response, err := us.repo.GetStates(ctx)
	return true, response, err
}
