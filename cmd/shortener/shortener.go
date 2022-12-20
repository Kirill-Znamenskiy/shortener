package main

import (
	"github.com/Kirill-Znamenskiy/shortener/internal/config"
	"github.com/Kirill-Znamenskiy/shortener/internal/handlers"
	"github.com/Kirill-Znamenskiy/shortener/internal/storage"
	"log"
	"net/http"
)

func main() {

	cfg := config.LoadEnvConfig()

	stg := storage.NewInMemoryStorage()

	mainHandler := handlers.MakeMainHandler(stg, cfg)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, mainHandler))
}
