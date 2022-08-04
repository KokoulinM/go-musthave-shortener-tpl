// Package server is a convenient wrapper over ListenAndServe
package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/handlers/middlewares"
)

type Server struct {
	// addr - contains the server address
	addr string
	// key - encryption key
	key []byte
	// handler - composable HTTP services with a large set of handlers.
	handler *chi.Mux
	// s- defines parameters for running an HTTP server.
	s *http.Server
}

// New is the server constructor
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

// Start is the method to start the server
func (s *Server) Start() error {
	err := s.s.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// Start is the method to start the server with tls
func (s *Server) StartTLS(certFile, keyFile string) error {
	err := s.s.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown is the method to stop the server
func (s *Server) Shutdown(ctx context.Context) error {
	err := s.s.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
