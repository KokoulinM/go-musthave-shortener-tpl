package server

import (
	"log"
	"net/http"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
)

type server struct {
	addr   string
	config configs.Config
}

func New(addr string, config configs.Config) *server {
	return &server{
		addr:   addr,
		config: config,
	}
}

func (s *server) Start() {
	h := handlers.New(s.config)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		router.Get("/{id}", h.Get)
		router.Get("/", h.Get)
		router.Post("/", h.Save)
		router.Post("/api/shorten", h.SaveJSON)
	})

	log.Fatal(http.ListenAndServe(s.addr, handlers.GzipHandle(router)))
}
