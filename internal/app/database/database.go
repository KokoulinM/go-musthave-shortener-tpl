package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/models"
)

type PostgresDatabase struct {
	conn    *sql.DB
	baseURL string
}

func NewRepository(baseURL string, db *sql.DB) *PostgresDatabase {
	return &PostgresDatabase{
		conn:    db,
		baseURL: baseURL,
	}
}

func New(baseURL string, db *sql.DB) handlers.Repository {
	return handlers.Repository(NewRepository(baseURL, db))
}

func (db *PostgresDatabase) AddURL(ctx context.Context, longURL models.LongURL, shortURL models.ShortURL, user models.UserID) error {
	return nil
}

func (db *PostgresDatabase) GetURL(ctx context.Context, shortURL models.ShortURL) (models.ShortURL, error) {
	return "", nil
}

func (db *PostgresDatabase) GetUserURLs(ctx context.Context, user models.UserID) ([]handlers.ResponseGetURL, error) {
	return nil, nil
}

func (db *PostgresDatabase) Ping(ctx context.Context) error {
	err := db.conn.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
