package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/workers"
)

func New(repo handlers.Repository, cfg configs.Config, wp *workers.WorkerPool) *chi.Mux {
	h := handlers.New(repo, cfg.BaseURL, wp)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		router.Post("/", h.CreateShortURL)
		router.Get("/{id}", h.RetrieveShortURL)
		router.Get("/ping", h.PingDB)
		router.Post("/api/shorten", h.ShortenURL)
		router.Get("/api/user/urls", h.GetUserURLs)
		router.Delete("/api/user/urls", h.DeleteBatch)
		router.Post("/api/shorten/batch", h.CreateBatch)
	})

	return router
}
