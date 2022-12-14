package main

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/internal/handlers"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/spf13/pflag"
	"log"
	"net/http"
)

func main() {

	cfg := config.LoadFromEnv()
	config.DefineFlags(cfg)
	pflag.Parse()

	var stg storage.Storage
	if cfg.StorageFilePath == "" {
		stg = storage.NewInMemoryStorage()
	} else {
		stg = storage.NewFileStorage(cfg.StorageFilePath)
	}

	mainHandler := handlers.MakeMainHandler(stg, cfg)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, mainHandler))
}
