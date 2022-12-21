package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"mime"
	"net/http"
	"net/url"
)

func MakeMainHandler(stg storage.Storage, cfg *config.Config) http.Handler {
	hs := &Handlers{stg: stg, cfg: cfg}
	r := chi.NewRouter()
	r.Use(middleware.ContentCharset("", "UTF-8"))
	r.Use(middleware.Compress(5, "text/html", "application/json"))
	r.Use(middleware.AllowContentType("", "text/plain", "text/html", "application/json"))
	r.Post("/", hs.makeWrapperForJSONHandlerFunc(hs.makeSaveNewURLHandlerFunc()))
	r.Post("/api/shorten", hs.makeSaveNewURLHandlerFunc())
	r.Get("/{key:[-\\w]+}", hs.makeGetURLHandlerFunc())
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})
	return r
}

type Handlers struct {
	stg storage.Storage
	cfg *config.Config
}

func (hs *Handlers) makeWrapperForJSONHandlerFunc(nextJSONHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		hContentType := req.Header.Get("Content-Type")

		isActive := false
		if hContentType == "" {
			isActive = true
		} else {
			ct, ctParams, err := mime.ParseMediaType(hContentType)
			fmt.Println("ctParams", ctParams)
			if err != nil || ct != "application/json" {
				isActive = true
			}
		}

		if isActive {
			req.Header.Set("Content-Type", "application/json")

			reqBodyCont, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			reqBodyCont = []byte(`{"URL":"` + string(reqBodyCont) + `"}`)
			req.Body = io.NopCloser(bytes.NewReader(reqBodyCont))

			crw := NewCustomResponseWriter()
			nextJSONHandler.ServeHTTP(crw, req)

			if crw.StatusCode == http.StatusCreated {
				respBodyData := new(struct {
					Result string `json:"result"`
				})
				err = json.Unmarshal(crw.Body, respBodyData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				crw.Body = []byte(respBodyData.Result)
				crw.Header().Set("Content-Type", "text/plain;charset=UTF-8")

				if respBodyData.Result == "" {
					http.Error(w, "Empty Result", http.StatusBadRequest)
					return
				}

				w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
				w.WriteHeader(http.StatusCreated)
			}

			_, err = crw.WriteToNext(w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			nextJSONHandler.ServeHTTP(w, req)
		}
	}
}

func (hs *Handlers) makeSaveNewURLHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		hContentType := req.Header.Get("Content-Type")
		if hContentType == "" {
			http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
			return
		}

		ct, ctParams, err := mime.ParseMediaType(hContentType)
		fmt.Println("ctParams", ctParams)
		if err != nil || ct != "application/json" {
			http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
			return
		}

		reqBodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reqBodyData := new(struct {
			URL string `json:"url"`
		})
		err = json.Unmarshal(reqBodyBytes, reqBodyData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reqBodyDataURL, err := url.Parse(reqBodyData.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key, err := hs.stg.Put("", reqBodyDataURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := hs.cfg.BaseURL + "/" + key

		respBodyData := new(struct {
			Result string `json:"result"`
		})
		respBodyData.Result = result

		respBodyBytes, err := json.Marshal(respBodyData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write(respBodyBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (hs *Handlers) makeGetURLHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//key := req.URL.EscapedPath()
		//key = strings.Trim(key, "/")
		key := chi.URLParam(req, "key")

		if toRespURL, isOk := hs.stg.Get(key); isOk {
			w.Header().Add("Location", toRespURL.String())
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "Resource Not Found", http.StatusBadRequest) // http.StatusNotFound
		}
	}
}
