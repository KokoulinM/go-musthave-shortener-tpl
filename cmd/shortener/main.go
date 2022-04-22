package main

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers/db"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	fmt.Println("main started")
	defer fmt.Println("main finished")

	cfg := configs.New()

	serv := server.New(cfg.ServerAddress, cfg)

	conn, err := db.Instance()
	if err != nil {
		log.Println("Closing connect to db")
		err := conn.Close()
		if err != nil {
			log.Println("Closing don't close")
		}
	}

	serv.Start()
}
