package main

import (
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	conf := configs.New()

	serv := server.New(conf.ServerAddress)

	serv.Start()
}
