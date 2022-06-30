package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers/middlewares"
)

type Server struct {
	addr    string
	key     []byte
	handler *chi.Mux
	s       *http.Server
}

func New(addr string, key []byte, handler *chi.Mux) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: middlewares.Conveyor(handler, middlewares.GzipMiddleware, middlewares.CookieMiddleware(key)),
	}

	return &Server{
		addr:    addr,
		key:     key,
		handler: handler,
		s:       srv,
	}
}

func (s *Server) Start() error {
	err := s.s.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.s.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
