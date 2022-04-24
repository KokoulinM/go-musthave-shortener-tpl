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
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/router"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	interrupt := make(chan os.Signal, 1)

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

	defer db.Conn.Close()

	handler := router.New(db, cfg)

	serv := server.New(db, cfg.ServerAddress, handler, cfg)

	go func() {
		serv.Start()
	}()

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		fmt.Println("Got SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("Got SIGTERM...")
	}
}
