package main

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/internal/handlers"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"log"
	"net/http"
)

func main() {

	cfg := config.LoadEnvConfig()

	var stg storage.Storage
	if cfg.StorageFilePath == "" {
		stg = storage.NewInMemoryStorage()
	} else {
		stg = storage.NewFileStorage(cfg.StorageFilePath)
	}

	mainHandler := handlers.MakeMainHandler(stg, cfg)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, mainHandler))
}
