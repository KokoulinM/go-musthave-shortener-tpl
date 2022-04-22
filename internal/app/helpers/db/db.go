package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
)

var instance *sql.DB

func Instance() (*sql.DB, error) {
	fmt.Println("db started")
	defer fmt.Println("db finished")

	if instance == nil {
		instance = new(sql.DB)

		dsn := configs.New()
		if dsn.DatabaseDSN == "" {
			return instance, fmt.Errorf("dsn can not be missing")
		}

		inst, err := sql.Open("pgx", dsn.DatabaseDSN)
		if err != nil {
			return instance, err
		}

		instance = inst

		log.Println("Connect to database")

		return instance, nil
	}

	return instance, nil
}
