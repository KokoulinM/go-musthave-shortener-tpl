package server

import (
	handler "github.com/KokoulinM/go-musthave-shortener-tpl/internal/handlers"
	"log"
	"net/http"
)

type Server struct {
	host    string
	handler handler.Handler
}

func New(host string) *Server {
	return &Server{
		host: host,
	}
}

func (s *Server) Start() {
	http.HandleFunc("/", s.handler.Save)
	http.HandleFunc("/{id:.+}", s.handler.Get)

	server := &http.Server{
		Addr: s.host,
	}

	log.Fatal(server.ListenAndServe())
}
