package main

import "github.com/KokoulinM/go-musthave-shortener-tpl/internal/server"

func main() {
	serv := server.New(":8080")

	serv.Start()
}
