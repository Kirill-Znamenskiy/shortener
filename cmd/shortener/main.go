package main

import (
	"fmt"
	"github.com/teris-io/shortid"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var shid2url map[string]*url.URL

func main() {

	shid2url = make(map[string]*url.URL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			reqBodyCont, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			reqBodyURL, err := url.Parse(string(reqBodyCont))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			shid, err := shortid.Generate()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			shid2url[shid] = reqBodyURL

			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			_, err = fmt.Fprint(w, "http://localhost:8080/"+shid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

		case http.MethodGet:
			shid := req.URL.EscapedPath()
			shid = strings.Trim(shid, "/")

			if toRespURL, isOk := shid2url[shid]; isOk {
				w.Header().Add("Location", toRespURL.String())
				w.WriteHeader(http.StatusTemporaryRedirect)
			} else {
				http.Error(w, "Resource Not Found", http.StatusBadRequest) // http.StatusNotFound
			}

		default:
			http.Error(w, "Method Not Allowed", http.StatusBadRequest) // http.StatusMethodNotAllowed
		}
	})

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
