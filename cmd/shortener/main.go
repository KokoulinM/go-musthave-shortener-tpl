package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/workers"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	"github.com/KokoulinM/go-musthave-shortener-tpl/cmd/shortener/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/cmd/shortener/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/cmd/shortener/router"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/storages"
)

func main() {
	var httpServer *server.Server

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	cfg := configs.New()

	var repo handlers.Repository

	wp := workers.New(ctx, cfg.Workers, cfg.WorkersBuffer)

	go func() {
		wp.Run(ctx)
	}()

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

	handler := router.New(repo, cfg, wp)

	g.Go(func() error {
		httpServer = server.New(cfg.ServerAddress, cfg.Key, handler)

		err := httpServer.Start()

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

	log.Println("Receive shutdown signal")

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if httpServer != nil {
		_ = httpServer.Shutdown()
	}

	err := g.Wait()
	if err != nil {
		log.Printf("server returning an error: %v", err)
		os.Exit(2)
	}
}
