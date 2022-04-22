package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers/db"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	fmt.Println("main started")
	defer fmt.Println("main finished")

	cfg := configs.New()

	conn, err := db.Conn("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Println("Closing connect to db")
		err := conn.Close()
		if err != nil {
			log.Println("Closing don't close")
		}
	}

	db := database.New(conn)

	serv := server.New(cfg.ServerAddress, cfg, db)

	go func() {
		serv.Start()
	}()

	interrupt := make(chan os.Signal, 1)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		fmt.Println("Got SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("Got SIGTERM...")
	}
}
