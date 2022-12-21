package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StorageFilePath string `env:"FILE_STORAGE_PATH"`
}

func LoadFromEnv() *Config {
	cfg := new(Config)
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
