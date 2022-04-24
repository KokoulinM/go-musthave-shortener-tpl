package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
)

func New(db *database.PostgresDatabase, cfg configs.Config) *chi.Mux {
	h := handlers.New(db, cfg)

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

	return router
}
