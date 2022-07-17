// Package configs is global configuration for application operation.
package configs

import (
	"flag"
	"log"
	"sync"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/helpers"
	"github.com/caarlos0/env/v6"
)

const (
	DefaultBaseURL         = "http://localhost:8080"
	DefaultServerAddress   = ":8080"
	DefaultFileStoragePath = "storage.json "
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

var (
	config Config
	once   sync.Once
)

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
	config = Config{
		baseURL:         DefaultBaseURL,
		serverAddress:   DefaultServerAddress,
		fileStoragePath: DefaultFileStoragePath,
		workers:         DefaultWorkers,
		workersBuffer:   DefaultWorkersBuffer,
	}

	return config
}

func New() *Config {
	once.Do(func() {
		config = defaultConfig()

		random, err := helpers.GenerateRandom(16)
		if err != nil {
			log.Fatal(err)
		}

		config.key = random

		err = env.Parse(&config)
		if err != nil {
			log.Fatal(err)
		}

		if checkExists("b") {
			flag.StringVar(&config.baseURL, "b", DefaultBaseURL, "BaseUrl")
		}

		if checkExists("a") {
			flag.StringVar(&config.serverAddress, "a", DefaultServerAddress, "ServerAddress")
		}

		if checkExists("f") {
			flag.StringVar(&config.fileStoragePath, "f", DefaultFileStoragePath, "FileStoragePath")
		}

		if checkExists("d") {
			flag.StringVar(&config.databaseDSN, "d", config.databaseDSN, "DatabaseDSN")
		}

		if checkExists("w") {
			flag.IntVar(&config.workers, "w", DefaultWorkers, "Workers")
		}

		if checkExists("wb") {
			flag.IntVar(&config.workersBuffer, "wb", DefaultWorkersBuffer, "WorkersBuffer")
		}

		flag.Parse()
	})

	return &config
}
