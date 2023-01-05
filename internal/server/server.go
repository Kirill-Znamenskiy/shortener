package server

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"log"
	"net/http"
)

func Run(cfg *config.Config) {

	mainHandler := MakeMainHandler(cfg)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, mainHandler))
}
