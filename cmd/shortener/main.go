package main

import (
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	serv := server.New("localhost", "8080")

	serv.Start()
}
