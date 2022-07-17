// Package configs is global configuration for application operation.
package configs

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/helpers"
)

const (
	DefaultBaseURL         = "http://localhost:8080"
	DefaultServerAddress   = ":8080"
	DefaultFileStoragePath = "storage.json "
	DefaultDatabaseDSN     = "user=postgres password=postgres sslmode=disable"
	DefaultWorkers         = 10
	DefaultWorkersBuffer   = 100
)

// Config contains app configuration.
type Config struct {
	// BaseURL - base app address
	baseURL string `env:"BASE_URL"`
	// ServerAddress - server address
	serverAddress string `env:"SERVER_ADDRESS"`
	// FileStoragePath - path to the file base
	fileStoragePath string `env:"FILE_STORAGE_PATH"`
	// DatabaseDSN - path to the database
	databaseDSN string `env:"DATABASE_DSN"`
	// Key - encryption key
	key []byte
	// Workers - number of workers
	workers int `env:"WORKERS"`
	// WorkersBuffer - buffer size value
	workersBuffer int `env:"WORKERS_BUFFER"`
}

func (c *Config) BaseURL() string {
	return c.baseURL
}

func (c *Config) ServerAddress() string {
	return c.serverAddress
}

func (c *Config) FileStoragePath() string {
	return c.fileStoragePath
}

func (c *Config) DatabaseDSN() string {
	return c.databaseDSN
}

func (c *Config) Key() []byte {
	return c.key
}

func (c *Config) Workers() int {
	return c.workers
}

func (c *Config) WorkersBuffer() int {
	return c.workersBuffer
}

// The function checks for the presence of a flag. f - flag values
func checkExists(f string) bool {
	return flag.Lookup(f) == nil
}

func defaultConfig() Config {
	return Config{
		baseURL:         DefaultBaseURL,
		serverAddress:   DefaultServerAddress,
		fileStoragePath: DefaultFileStoragePath,
		databaseDSN:     DefaultDatabaseDSN,
		workers:         DefaultWorkers,
		workersBuffer:   DefaultWorkersBuffer,
	}
}

func New() Config {
	c := defaultConfig()

	random, err := helpers.GenerateRandom(16)
	if err != nil {
		log.Fatal(err)
	}

	c.key = random

	err = env.Parse(&c)
	if err != nil {
		log.Fatal(err)
	}

	if checkExists("b") {
		flag.StringVar(&c.baseURL, "b", DefaultBaseURL, "BaseUrl")
	}

	if checkExists("a") {
		flag.StringVar(&c.serverAddress, "a", DefaultServerAddress, "ServerAddress")
	}

	if checkExists("f") {
		flag.StringVar(&c.fileStoragePath, "f", DefaultFileStoragePath, "FileStoragePath")
	}

	if checkExists("d") {
		flag.StringVar(&c.databaseDSN, "d", DefaultDatabaseDSN, "DatabaseDSN")
	}

	if checkExists("w") {
		flag.IntVar(&c.workers, "w", DefaultWorkers, "Workers")
	}

	if checkExists("wb") {
		flag.IntVar(&c.workersBuffer, "wb", DefaultWorkersBuffer, "WorkersBuffer")
	}

	flag.Parse()

	return c
}
