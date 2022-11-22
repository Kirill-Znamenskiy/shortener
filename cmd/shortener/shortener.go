package main

import (
	"github.com/Kirill-Znamenskiy/shortener/internal/handlers"
	"github.com/Kirill-Znamenskiy/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	stg := storage.NewInMemoryStorage()

	rootHandler := handlers.MakeRootHandler(stg)

	r := chi.NewRouter()
	r.Handle("/*", rootHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
