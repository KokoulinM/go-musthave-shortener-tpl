package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	_ "net/http/pprof"

	_ "github.com/lib/pq"

	"golang.org/x/sync/errgroup"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/database/filebase"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/database/postgres"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/router"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/server"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/workers"
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
	defer wp.Stop()
	if cfg.DatabaseDSN != "" {
		conn, err := postgres.Conn("postgres", cfg.DatabaseDSN)
		if err != nil {
			log.Printf("Unable to connect to the database: %s", err.Error())
		}

		err = postgres.SetUpDataBase(ctx, conn)

		if err != nil {
			log.Printf("Unable to create database struct: %s", err.Error())
		}

		repo = postgres.NewDatabaseRepository(cfg.BaseURL, conn)
	} else {
		repo = filebase.NewFileRepository(ctx, cfg.FileStoragePath, cfg.BaseURL)
	}

	g, ctx := errgroup.WithContext(ctx)

	h := handlers.New(repo, cfg.BaseURL, wp)

	mux := router.New(h)

	g.Go(func() error {
		httpServer = server.New(cfg.ServerAddress, cfg.Key, mux)

		err := httpServer.Start()
		if err != nil {
			return err
		}

		log.Printf("httpServer starting at: %v", cfg.ServerAddress)

		return nil
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Println("Receive shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if httpServer != nil {
		_ = httpServer.Shutdown(shutdownCtx)
	}

	sl := []string{"foo", "bar", "buzz"}
	sl = sort.StringSlice(sl) // sort.StringSlice — это не функция, а тип, выражение не отсортирует sl
	// чтобы отсортировать, нужно сделать sort.StringSlice(sl).Sort()

	err := g.Wait()
	if err != nil {
		log.Printf("server returning an error: %v", err)
		os.Exit(2)
	}
}
