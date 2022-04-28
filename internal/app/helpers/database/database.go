package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func Conn(driverName, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn can not be missing")
	}

	if driverName == "" {
		return nil, fmt.Errorf("driver name can not be missing")
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return db, err
	}

	log.Println("Connect to database")

	return db, nil
}

func SetUpDataBase(db *sql.DB, ctx context.Context) error {
	//var extention string
	//
	//query := db.QueryRowContext(ctx, "SELECT 'exists' FROM pg_extension WHERE extname='uuid-ossp';")
	//
	//err := query.Scan(&extention)
	//if err != nil {
	//	log.Printf("Unable to scan: %s", err)
	//}
	//
	//if extention != "exists" {
	//	_, err := db.ExecContext(ctx, `CREATE EXTENSION "uuid-ossp";`)
	//	if err != nil {
	//		return err
	//	}
	//
	//	log.Println("Create EXTENSION")
	//}

	sqlCreateDB := `CREATE TABLE IF NOT EXISTS urls (
								id serial PRIMARY KEY,
								user_id uuid DEFAULT uuid_generate_v4 (), 	
								origin_url VARCHAR NOT NULL, 
								short_url VARCHAR NOT NULL,
								is_deleted BOOLEAN NOT NULL DEFAULT FALSE
					);`
	res, err := db.ExecContext(ctx, sqlCreateDB)

	log.Println("Create table", err, res)

	return nil
}
