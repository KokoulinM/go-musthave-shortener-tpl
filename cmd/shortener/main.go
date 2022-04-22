package main

import (
	"fmt"
	"log"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers/db"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	fmt.Println("main started")
	defer fmt.Println("main finished")

	c := configs.New()

	serv := server.New(c.ServerAddress, c)

	conn, err := db.New(c.DatabaseDSN)
	if err != nil {
		log.Println("Closing connect to db")
		err := conn.Close()
		if err != nil {
			log.Println("Closing don't close")
		}
	}

	serv.Start()
}
