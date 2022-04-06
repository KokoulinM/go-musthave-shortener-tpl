package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
)

type server struct{}

func New() *server {
	return &server{}
}

func (s *server) Start() {
	handlers := handlers.New()

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		router.Get("/{id}", handlers.Get)
		router.Get("/", handlers.Get)
		router.Post("/", handlers.Save)
		router.Post("/api/shorten", handlers.SaveJSON)
	})

	log.Fatal(http.ListenAndServe(handlers.Config.GetServerAddress(), router))
}
