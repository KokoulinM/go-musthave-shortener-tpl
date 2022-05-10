package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers/middlewares"
)

type server struct {
	addr    string
	key     []byte
	handler *chi.Mux
}

func New(addr string, key []byte, handler *chi.Mux) *server {
	return &server{
		addr:    addr,
		key:     key,
		handler: handler,
	}
}

func (s *server) Start() error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: middlewares.Conveyor(s.handler, middlewares.GzipMiddleware, middlewares.CookieMiddleware(s.key)),
	}

	if err := http.ListenAndServe(srv.Addr, srv.Handler); err != http.ErrServerClosed {
		return err
	}

	return nil
}
