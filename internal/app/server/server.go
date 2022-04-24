package server

import (
	"fmt"
	"log"
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
	fmt.Println("server started")
	defer fmt.Println("server finished")

	return &server{
		addr:    addr,
		key:     key,
		handler: handler,
	}
}

func (s *server) Start() {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: middlewares.Conveyor(s.handler, middlewares.GzipMiddleware, middlewares.CookieMiddleware(s.key)),
	}

	log.Fatal(http.ListenAndServe(srv.Addr, srv.Handler))
}
