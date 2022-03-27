package server

import (
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers"
	"log"
	"net/http"
)

type Server struct {
	host string
}

func New(host string) *Server {
	return &Server{
		host: host,
	}
}

func (s *Server) Start() {
	handler := handlers.New()

	http.HandleFunc("/", handler.CommonHandler)

	server := &http.Server{
		Addr: s.host,
	}

	log.Fatal(server.ListenAndServe())
}
