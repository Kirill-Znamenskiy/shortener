package main

import (
	"github.com/Kirill-Znamenskiy/shortener/internal/handlers"
	"github.com/Kirill-Znamenskiy/shortener/internal/storage"
	"log"
	"net/http"
)

func main() {

	stg := storage.NewInMemoryStorage()
	hs := handlers.Handlers{Stg: stg}

	log.Fatal(http.ListenAndServe("localhost:8080", hs.MakeMainHandler()))
}
