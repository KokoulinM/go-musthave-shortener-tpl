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

	if driverName == "" {
		return instance, fmt.Errorf("driver name can not be missing")
	}

	instance, err := sql.Open(driverName, dsn)
	if err != nil {
		return instance, err
	}

	log.Println("Connect to database")

	return instance, nil
}
