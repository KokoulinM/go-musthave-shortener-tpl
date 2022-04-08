package configs

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func New() Config {
	var c Config

	err := env.Parse(&c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
