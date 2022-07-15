// Package router composable HTTP services with a large set of handlers
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/handlers"
)

// New router constructor
func New(h *handlers.Handlers) *chi.Mux {
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
