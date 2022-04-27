package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/storages"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/router"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	cfg := configs.New()

	var repo handlers.Repository

	if cfg.DatabaseDSN != "" {
		conn, err := database.Conn("pgx", cfg.DatabaseDSN)
		if err != nil {
			log.Println("Closing don't close")
		}

		database.SetUpDataBase(conn, ctx)

		repo = storages.NewDatabaseRepository(cfg.DatabaseDSN, conn)
	} else {
		repo = storages.NewFileRepository(ctx, cfg.FileStoragePath, cfg.BaseURL)
	}

	handler := router.New(repo, cfg)

	serv := server.New(cfg.ServerAddress, cfg.Key, handler)

	go func() {
		serv.Start()
	}()

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}
}
