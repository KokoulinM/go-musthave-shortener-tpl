package db

import (
	"database/sql"
	"fmt"
	"log"
)

func Conn(driverName, dsn string) (*sql.DB, error) {
	instance := new(sql.DB)

	if dsn == "" {
		return instance, fmt.Errorf("dsn can not be missing")
	}

	inst, err := sql.Open(driverName, dsn)
	if err != nil {
		return instance, err
	}

	instance = inst

	log.Println("Connect to database")

	return instance, nil
}
