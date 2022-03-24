package server

import (
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
	http.HandleFunc("/", Handler)
	http.HandleFunc("/{id:.+}", Handler)

	server := http.Server{
		Addr: s.host,
	}

	log.Fatal(server.ListenAndServe())
}
