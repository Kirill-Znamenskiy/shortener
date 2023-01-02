package main

import (
	"context"
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/internal/server"
	"log"
)

func main() {

	ctx := context.Background()

	cfg := new(config.Config)
	config.LoadFromEnv(ctx, cfg)
	config.ParseFlags(cfg)

	log.Printf("Now run API server with config: %s\n", config.ToPrettyString(cfg))

	server.Run(cfg)
}
