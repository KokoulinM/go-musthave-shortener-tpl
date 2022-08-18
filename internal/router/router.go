// Package router composable HTTP services with a large set of handlers
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/configs"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/handlers"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/handlers/middlewares"
)

// New router constructor
func New(h *handlers.Handlers, cfg *configs.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Post("/", h.CreateShortURL)
		r.Get("/{id}", h.RetrieveShortURL)
		r.Get("/ping", h.PingDB)
		r.Post("/api/shorten", h.ShortenURL)
		r.Get("/api/user/urls", h.GetUserURLs)
		r.Delete("/api/user/urls", h.DeleteBatch)
		r.Post("/api/shorten/batch", h.CreateBatch)
		router.Route("/api/internal", func(r chi.Router) {
			r.Use(middlewares.InternalMiddleware(cfg.TrustedSubnet))
			r.Get("/stats", h.GetStates)
		})

	})

	return router
}
