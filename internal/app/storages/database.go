package storages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/models"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/shortener"
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
	sqlAddRow := `INSERT INTO urls (user_id, origin_url, short_url)
				  VALUES ($1, $2, $3)`

	_, err := db.conn.ExecContext(ctx, sqlAddRow, user, longURL, shortURL)

	if err.(pgx.PgError).Code == pgerrcode.UniqueViolation {
		return handlers.NewErrorWithDB(err, "UniqConstraint")
	}

	return err
}

func (db *PostgresDatabase) AddMultipleURLs(ctx context.Context, urls []handlers.RequestGetURLs, user models.UserID) ([]handlers.ResponseGetURLs, error) {
	var result []handlers.ResponseGetURLs
	tx, err := db.conn.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO urls (user_id, origin_url, short_url) VALUES ($1, $2, $3)`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for _, u := range urls {
		shortURL := shortener.ShorterURL(u.OriginalURL)
		if _, err = stmt.ExecContext(ctx, user, u.OriginalURL, shortURL); err != nil {
			return nil, err
		}
		result = append(result, handlers.ResponseGetURLs{
			CorrelationID: u.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", db.baseURL, shortURL),
		})
	}

	if err != nil {
		return nil, err
	}
	tx.Commit()
	return result, nil
}

func (db *PostgresDatabase) GetURL(ctx context.Context, shortURL models.ShortURL) (models.ShortURL, error) {
	sqlGetURLRow := `SELECT origin_url FROM urls WHERE short_url=$1 LIMIT 1`

	row := db.conn.QueryRowContext(ctx, sqlGetURLRow, shortURL)

	result := GetURLData{}

	err := row.Scan(&result.OriginalURL)
	if err != nil {
		return "", err
	}

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

	sqlGetUserURL := `SELECT origin_url, short_url FROM urls WHERE user_id=$1;`
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
