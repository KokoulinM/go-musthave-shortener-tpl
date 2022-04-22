package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers/middlewares"
)

type server struct {
	addr   string
	config configs.Config
	db     *database.PostgresDatabase
}

func New(addr string, config configs.Config, db *database.PostgresDatabase) *server {
	fmt.Println("server started")
	defer fmt.Println("server finished")

	return &server{
		addr:   addr,
		config: config,
		db:     db,
	}
}

func (s *server) Start() {
	h := handlers.New(s.config, s.db)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		router.Get("/{id}", h.Get)
		router.Get("/", h.Get)
		router.Post("/", h.Save)
		router.Post("/api/shorten", h.SaveJSON)
		router.Get("/api/user/urls", h.GetLinks)
		router.Get("/ping", h.PingDB)
	})

	log.Fatal(http.ListenAndServe(s.addr, middlewares.Conveyor(router, middlewares.GzipMiddleware, middlewares.CookieMiddleware)))
}
