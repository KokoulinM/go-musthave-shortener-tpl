package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/filebase"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers/db"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/router"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	cfg := configs.New()

	var repo handlers.Repository

	log.Println(cfg)

	if cfg.DatabaseDSN != "" {
		conn, err := db.Conn("pgx", cfg.DatabaseDSN)
		if err != nil {
			log.Println("Closing connect to db")
			err := conn.Close()
			if err != nil {
				log.Println("Closing don't close")
			}
		}
		defer conn.Close()

		repo = database.New(cfg.DatabaseDSN, conn)
	} else {
		repo = filebase.New(ctx, cfg.FileStoragePath, cfg.BaseURL)
	}

	handler := router.New(repo, cfg)

	serv := server.New(cfg.ServerAddress, cfg.Key, handler)

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
