package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KokoulinM/go-musthave-shortener-tpl/cmd/shortener/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/workers"
)

func New(repo handlers.Repository, cfg configs.Config, wp *workers.WorkerPool) *chi.Mux {
	h := handlers.New(repo, cfg.BaseURL, wp)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		router.Get("/", h.Get)
		router.Post("/", h.Save)
		router.Get("/{id}", h.Get)
		router.Get("/ping", h.PingDB)
		router.Post("/api/shorten", h.SaveJSON)
		router.Get("/api/user/urls", h.GetLinks)
		router.Delete("/api/user/urls", h.DeleteLinks)
		router.Post("/api/shorten/batch", h.CreateBatch)
	})

	return router
}
