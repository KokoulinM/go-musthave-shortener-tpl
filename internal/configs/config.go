package configs

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/helpers"
)

type Config struct {
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage.json"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             []byte
	Workers         int `env:"WORKERS" envDefault:"10"`
	WorkersBuffer   int `env:"WORKERS_BUFFER" envDefault:"100"`
}

func checkExists(f string) bool {
	return flag.Lookup(f) == nil
}

func New() Config {
	var c Config

	random, err := helpers.GenerateRandom(16)
	if err != nil {
		log.Fatal(err)
	}

	c.Key = random

	err = env.Parse(&c)
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

	if checkExists("w") {
		flag.IntVar(&c.Workers, "w", c.Workers, "Workers")
	}

	if checkExists("wb") {
		flag.IntVar(&c.WorkersBuffer, "wb", c.WorkersBuffer, "WorkersBuffer")
	}

	flag.Parse()

	return c
}
