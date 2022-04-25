package configs

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers"
)

type Config struct {
	BaseURL         string `env:"BASE_URL" envDefault:":8080"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage.json"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             []byte
}

func checkExists(f string) bool {
	return flag.Lookup(f) == nil
}

func New() Config {
	var c Config

	c.Key, _ = helpers.GenerateRandom(16)

	err := env.Parse(&c)
	if err != nil {
		log.Fatal(err)
	}

	if checkExists("b") {
		flag.StringVar(&c.BaseURL, "b", c.BaseURL, "BaseUrl")
	}

	if checkExists("a") {
		flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "ServerAddress")
	}

	if checkExists("f") {
		flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "FileStoragePath")
	}

	if checkExists("d") {
		flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "DatabaseDSN")
	}

	flag.Parse()

	return c
}
