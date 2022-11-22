package main

import (
	"github.com/Kirill-Znamenskiy/shortener/internal/handlers"
	"github.com/Kirill-Znamenskiy/shortener/internal/storage"
	"log"
	"net/http"
)

func main() {

	stg := storage.NewInMemoryStorage()

	rootHandler := handlers.MakeRootHandler(stg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
