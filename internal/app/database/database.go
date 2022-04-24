package database

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresDatabase struct {
	Conn *sql.DB
}

func New(conn *sql.DB) *PostgresDatabase {
	return &PostgresDatabase{
		Conn: conn,
	}
}

func (db *PostgresDatabase) Ping(ctx context.Context) error {
	err := db.Conn.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
