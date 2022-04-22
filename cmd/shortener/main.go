package main

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	fmt.Println("main started")
	defer fmt.Println("main finished")

	cfg := configs.New()

	serv := server.New(cfg.ServerAddress, cfg)

	serv.Start()
}
