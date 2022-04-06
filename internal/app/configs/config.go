package configs

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config interface {
	GetBaseURL() string
	GetServerAddress() string
}

type config struct {
	baseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	serverAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
}

func (c config) GetBaseURL() string {
	return c.baseURL
}

func (c config) GetServerAddress() string {
	return c.serverAddress
}

func New() config {
	var cfg config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
