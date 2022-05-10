package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers/middlewares"
)

type Server struct {
	addr    string
	key     []byte
	handler *chi.Mux
}

func New(addr string, key []byte, handler *chi.Mux) *Server {
	return &Server{
		addr:    addr,
		key:     key,
		handler: handler,
	}
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: middlewares.Conveyor(s.handler, middlewares.GzipMiddleware, middlewares.CookieMiddleware(s.key)),
	}

	if err := http.ListenAndServe(srv.Addr, srv.Handler); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	err := s.Shutdown()
	if err != nil {
		return err
	}

	return nil
}
