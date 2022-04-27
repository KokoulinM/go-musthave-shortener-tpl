package storages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/models"
)

type PostgresDatabase struct {
	conn    *sql.DB
	baseURL string
}

type GetURLData struct {
	OriginalURL string
	IsDeleted   bool
}

func DatabaseRepository(baseURL string, db *sql.DB) *PostgresDatabase {
	return &PostgresDatabase{
		conn:    db,
		baseURL: baseURL,
	}
}

func NewDatabaseRepository(baseURL string, db *sql.DB) handlers.Repository {
	return handlers.Repository(DatabaseRepository(baseURL, db))
}

func (db *PostgresDatabase) AddURL(ctx context.Context, longURL models.LongURL, shortURL models.ShortURL, user models.UserID) error {
	sqlAddRow := `INSERT INTO urls (user_id, original_url, short_url)
				  VALUES ($1, $2, $3)`

	_, err := db.conn.ExecContext(ctx, sqlAddRow, user, longURL, shortURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add an entry to the table: %v\n", err)
	}

	return err
}

func (db *PostgresDatabase) GetURL(ctx context.Context, shortURL models.ShortURL) (models.ShortURL, error) {
	sqlGetURLRow := `SELECT original_url, is_deleted FROM urls WHERE short_url=$1 LIMIT 1`

	row := db.conn.QueryRowContext(ctx, sqlGetURLRow, shortURL)

	result := GetURLData{}

	row.Scan(&result.OriginalURL, &result.IsDeleted)

	if result.OriginalURL == "" {
		return "", errors.New("not found")
	}
	if result.IsDeleted {
		return "", errors.New("deleted")
	}

	return result.OriginalURL, nil
}

func (db *PostgresDatabase) GetUserURLs(ctx context.Context, user models.UserID) ([]handlers.ResponseGetURL, error) {
	var result []handlers.ResponseGetURL

	sqlGetUserURL := `SELECT original_url, short_url FROM urls WHERE user_id=$1 AND is_deleted=false;`
	rows, err := db.conn.QueryContext(ctx, sqlGetUserURL, user)
	if err != nil {
		return result, err
	}
	if rows.Err() != nil {
		return result, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var u handlers.ResponseGetURL
		err = rows.Scan(&u.OriginalURL, &u.ShortURL)
		if err != nil {
			return result, err
		}
		u.ShortURL = db.baseURL + u.ShortURL
		result = append(result, u)
	}

	return result, nil
}

func (db *PostgresDatabase) Ping(ctx context.Context) error {
	err := db.conn.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
