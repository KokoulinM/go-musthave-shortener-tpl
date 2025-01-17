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
	DefaultTrustedSubnet   = "127.0.0.1/24"
	DefaultGRPCPort        = 5000
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
	WorkersBuffer int  `env:"WORKERS_BUFFER"`
	EnableHttps   bool `env:"ENABLE_HTTPS" json:"ENABLE_HTTPS"`
	// Config - configuration file
	Config string `env:"CONFIG"`
	// TrustedSubnet - available url for internal requests
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	GRPCPort      int    `env:"GRPC_PORT" json:"grpc_port"`
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
		TrustedSubnet:   DefaultTrustedSubnet,
		GRPCPort:        DefaultGRPCPort,
	}
}

func readCfgFile(name string, cfg *Config) error {
	jsonFile, err := os.Open(name)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return err
	}

	return nil
}

func New() *Config {
	c := defaultConfig()

	random, err := helpers.GenerateRandom(16)
	if err != nil {
		log.Fatal(err)
	}

	c.Key = random

	if checkExists("s") {
		flag.StringVar(&c.Config, "s", c.Config, "Config")
	}

	cfgPath := os.Getenv("CONFIG")

	if cfgPath != "" {
		c.Config = cfgPath
	}

	if c.Config != "" {
		err := readCfgFile(c.Config, &c)
		if err != nil {
			log.Fatal("An error occurred while reading the configuration file")
		}
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

	if checkExists("t") {
		flag.StringVar(&c.TrustedSubnet, "t", c.TrustedSubnet, "TrustedSubnet")
	}

	flag.Parse()

	return &c
}
