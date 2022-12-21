package main

import (
	"flag"
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/internal/handlers"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"log"
	"net/http"
)

func main() {

	cfg := config.LoadFromEnv()

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.ServerAddress, "server-address", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url")
	flag.StringVar(&cfg.BaseURL, "base-url", cfg.BaseURL, "base url")
	flag.StringVar(&cfg.StorageFilePath, "f", cfg.StorageFilePath, "storage file path")
	flag.StringVar(&cfg.StorageFilePath, "storage-file-path", cfg.StorageFilePath, "storage file path")

	flag.Parse()

	var stg storage.Storage
	if cfg.StorageFilePath == "" {
		stg = storage.NewInMemoryStorage()
	} else {
		stg = storage.NewFileStorage(cfg.StorageFilePath)
	}

	mainHandler := handlers.MakeMainHandler(stg, cfg)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, mainHandler))
}
