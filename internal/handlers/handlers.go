package handlers

import (
	"fmt"
	"github.com/Kirill-Znamenskiy/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/teris-io/shortid"
	"io"
	"net/http"
	"net/url"
)

func MakeMainHandler(stg storage.Storage) http.Handler {
	hs := HandlerFuncs{stg: stg}
	r := chi.NewRouter()
	r.Post("/", hs.SaveNewURLHandlerFunc)
	r.Get("/{shid:[-\\w]+}", hs.GetURLHandlerFunc)
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})
	return r
}

type HandlerFuncs struct {
	stg storage.Storage
}

func (hs HandlerFuncs) SaveNewURLHandlerFunc(w http.ResponseWriter, req *http.Request) {
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
		_, isOk := hs.stg.Get(shid)
		if !isOk {
			break
		}
	}
	hs.stg.Put(shid, reqBodyURL)

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprint(w, "http://localhost:8080/"+shid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (hs HandlerFuncs) GetURLHandlerFunc(w http.ResponseWriter, req *http.Request) {
	//shid := req.URL.EscapedPath()
	//shid = strings.Trim(shid, "/")
	shid := chi.URLParam(req, "shid")

	if toRespURL, isOk := hs.stg.Get(shid); isOk {
		w.Header().Add("Location", toRespURL.String())
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Resource Not Found", http.StatusBadRequest) // http.StatusNotFound
	}
}
