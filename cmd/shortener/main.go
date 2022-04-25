package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/filebase"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers/db"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/router"
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
		conn, err := db.Conn("pgx", cfg.DatabaseDSN)
		if err != nil {
			log.Println("Closing don't close")
		}

		repo = database.New(cfg.DatabaseDSN, conn)
	} else {
		repo = filebase.New(ctx, cfg.FileStoragePath, cfg.BaseURL)
	}

	handler := router.New(repo, cfg)

	//serv := server.New(cfg.ServerAddress, cfg.Key, handler)

	go func() error {
		httpServer := &http.Server{
			Addr:    cfg.ServerAddress,
			Handler: handler,
		}
		log.Printf("httpServer starting at: %v", cfg.ServerAddress)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	}()

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}
}
