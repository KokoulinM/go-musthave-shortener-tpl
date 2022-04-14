package main

import (
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	c := configs.New()

	serv := server.New(c.ServerAddress, c)

	serv.Start()
}
