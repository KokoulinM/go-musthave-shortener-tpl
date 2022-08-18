// Package configs is global configuration for application operation.
package configs

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/caarlos0/env/v6"

	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/helpers"
)

const (
	DefaultBaseURL         = "http://localhost:8080"
	DefaultServerAddress   = ":8080"
	DefaultFileStoragePath = "storage.json "
	DefaultWorkers         = 10
	DefaultWorkersBuffer   = 100
	DefaultEnableHttps     = false
)

// Config contains app configuration.
type Config struct {
	// BaseURL - base app address
	BaseURL string `env:"BASE_URL" json:"BASE_URL"`
	// ServerAddress - server address
	ServerAddress string `env:"SERVER_ADDRESS" json:"SERVER_ADDRESS"`
	// FileStoragePath - path to the file base
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"FILE_STORAGE_PATH"`
	// DatabaseDSN - path to the database
	DatabaseDSN string `env:"DATABASE_DSN" json:"DATABASE_DSN"`
	// Key - encryption key
	Key []byte
	// Workers - number of workers
	Workers int `env:"WORKERS"`
	// WorkersBuffer - buffer size value
	WorkersBuffer int    `env:"WORKERS_BUFFER"`
	EnableHttps   bool   `env:"ENABLE_HTTPS" json:"ENABLE_HTTPS"`
	Config        string `env:"CONFIG"`
}

// The function checks for the presence of a flag. f - flag values
func checkExists(f string) bool {
	return flag.Lookup(f) == nil
}

func defaultConfig() Config {
	return Config{
		BaseURL:         DefaultBaseURL,
		ServerAddress:   DefaultServerAddress,
		FileStoragePath: DefaultFileStoragePath,
		Workers:         DefaultWorkers,
		WorkersBuffer:   DefaultWorkersBuffer,
		EnableHttps:     DefaultEnableHttps,
	}
}

func readCfgFile(name string, cfg *Config) {
	jsonFile, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func New() *Config {
	c := defaultConfig()

	random, err := helpers.GenerateRandom(16)
	if err != nil {
		log.Fatal(err)
	}

	c.Key = random

	if checkExists("c") {
		flag.StringVar(&c.Config, "s", c.Config, "Config")
	}

	cfgPath := os.Getenv("CONFIG")

	if cfgPath != "" {
		c.Config = cfgPath
	}

	if c.Config != "" {
		readCfgFile(c.Config, &c)
	}

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

	if checkExists("s") {
		flag.BoolVar(&c.EnableHttps, "s", c.EnableHttps, "EnableHttps")
	}

	flag.Parse()

	return &c
}
