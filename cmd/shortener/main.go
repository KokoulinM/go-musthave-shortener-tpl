package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	"github.com/KokoulinM/go-musthave-shortener-tpl/cmd/shortener/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/cmd/shortener/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/cmd/shortener/router"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/storages"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/workers"
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
		conn, err := database.Conn("postgres", cfg.DatabaseDSN)
		if err != nil {
			log.Printf("Unable to connect to the database: %s", err.Error())
		}

		err = database.SetUpDataBase(conn, ctx)

		if err != nil {
			log.Printf("Unable to create database struct: %s", err.Error())
		}

		repo = storages.NewDatabaseRepository(cfg.BaseURL, conn)
	} else {
		repo = storages.NewFileRepository(ctx, cfg.FileStoragePath, cfg.BaseURL)
	}

	g, ctx := errgroup.WithContext(ctx)

	wp := workers.New(ctx, cfg.Workers, cfg.WorkersBuffer)

	handler := router.New(repo, cfg, *wp)

	g.Go(func() error {
		serv := server.New(cfg.ServerAddress, cfg.Key, handler)

		err := serv.Start()

		log.Printf("httpServer starting at: %v", cfg.ServerAddress)

		if err != nil {
			return err
		}

		return nil
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	err := g.Wait()
	if err != nil {
		log.Printf("server returning an error: %v", err)
		os.Exit(2)
	}
}
