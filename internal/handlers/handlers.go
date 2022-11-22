package handlers

import (
	"fmt"
	"github.com/Kirill-Znamenskiy/shortener/internal/storage"
	"github.com/teris-io/shortid"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func RootHandler(w http.ResponseWriter, req *http.Request, stg storage.Storage) {
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

		shid := ""
		for {
			shid, err = shortid.Generate()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, isOk := stg.Get(shid)
			if !isOk {
				break
			}
		}
		stg.Put(shid, reqBodyURL)

		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprint(w, "http://localhost:8080/"+shid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

	case http.MethodGet:
		shid := req.URL.EscapedPath()
		shid = strings.Trim(shid, "/")

		if toRespURL, isOk := stg.Get(shid); isOk {
			w.Header().Add("Location", toRespURL.String())
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "Resource Not Found", http.StatusBadRequest) // http.StatusNotFound
		}

	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest) // http.StatusMethodNotAllowed
	}
}

func MakeRootHandler(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RootHandler(w, r, s)
	}
}
