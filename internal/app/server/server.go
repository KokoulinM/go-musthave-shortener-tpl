package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers/middlewares"
)

type server struct {
	addr    string
	config  configs.Config
	db      *database.PostgresDatabase
	handler *chi.Mux
}

func New(db *database.PostgresDatabase, addr string, handler *chi.Mux, config configs.Config) *server {
	fmt.Println("server started")
	defer fmt.Println("server finished")

	return &server{
		addr:    addr,
		config:  config,
		db:      db,
		handler: handler,
	}
}

func (s *server) Start() {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: middlewares.Conveyor(s.handler, middlewares.GzipMiddleware, middlewares.CookieMiddleware),
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if srv != nil {
		_ = srv.Shutdown(shutdownCtx)
	}

	log.Fatal(http.ListenAndServe(srv.Addr, srv.Handler))
}
