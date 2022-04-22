package database

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresDatabase struct {
	conn *sql.DB
}

func New(conn *sql.DB) *PostgresDatabase {
	return &PostgresDatabase{
		conn: conn,
	}
}

func (db *PostgresDatabase) Ping(ctx context.Context) error {
	err := db.conn.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
