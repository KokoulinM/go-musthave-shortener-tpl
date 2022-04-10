package configs

import (
	"flag"
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

	if flag.Lookup("b") == nil {
		flag.StringVar(&c.BaseURL, "b", c.BaseURL, "BaseUrl")
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "ServerAddress")
	}

	if flag.Lookup("f") == nil {
		flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "FileStoragePath")
	}

	flag.Parse()

	return c
}
