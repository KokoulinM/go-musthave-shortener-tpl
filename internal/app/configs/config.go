package configs

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
}

func New() Config {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
