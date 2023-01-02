package server

import (
	"bytes"
	"encoding/json"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic"
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"mime"
	"net/http"
)

const (
	UTF8                       = "UTF-8"
	TextHTML                   = "text/html"
	TextPlain                  = "text/plain"
	ApplicationJSON            = "application/json"
	TextPlainCharsetUTF8       = "text/plain;charset=UTF-8"
	ApplicationJSONCharsetUTF8 = "application/json;charset=UTF-8"
)

func MakeMainHandler(cfg *config.Config) http.Handler {

	hs := &Handlers{cfg: cfg, stg: cfg.GetStorage()}

	r := chi.NewRouter()

	r.Use(CleanURLPathMiddleware())
	r.Use(AllowContentCharsetMiddleware("", UTF8))

	r.Use(DecompressMiddleware())
	r.Use(middleware.Compress(5, TextHTML, TextPlain, ApplicationJSON))

	r.Use(middleware.AllowContentType("", TextHTML, TextPlain, ApplicationJSON))

	r.Post("/", hs.makeWrapperForJSONHandlerFunc(hs.makeSaveNewURLHandlerFunc()))
	r.Post("/api/shorten", hs.makeSaveNewURLHandlerFunc())
	r.Get("/{key:[-\\w]+}", hs.makeGetURLHandlerFunc())

	r.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("UNEXPECTED: %s %q\n", req.Method, req.URL.Path)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	return r
}

type Handlers struct {
	cfg *config.Config
	stg storage.Storage
}

func (hs *Handlers) makeWrapperForJSONHandlerFunc(nextJSONHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		hContentType := req.Header.Get("Content-Type")

		isActive := false
		if hContentType == "" {
			isActive = true
		} else {
			ct, _, err := mime.ParseMediaType(hContentType)
			if err != nil || ct != ApplicationJSON {
				isActive = true
			}
		}

		if isActive {

			reqBodyCont, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			reqBodyCont = []byte(`{"URL":"` + string(reqBodyCont) + `"}`)
			req.Body = io.NopCloser(bytes.NewReader(reqBodyCont))

			req.Header.Set("Content-Type", ApplicationJSON)

			crw := NewCustomResponseWriter()
			nextJSONHandler.ServeHTTP(crw, req)

			if crw.StatusCode == http.StatusCreated {
				respData := new(struct {
					Result string `json:"result"`
				})
				err = json.Unmarshal(crw.Body, respData)
				if checkErrorAsInternalServerError(w, err) {
					return
				}

				crw.Body = []byte(respData.Result)
				crw.Header().Set("Content-Type", TextPlainCharsetUTF8)

				if respData.Result == "" {
					http.Error(w, "Empty Result", http.StatusBadRequest)
					return
				}
			}

			_, err = crw.WriteToNext(w)
			if checkErrorAsInternalServerError(w, err) {
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

		ct, _, err := mime.ParseMediaType(hContentType)
		if err != nil || ct != ApplicationJSON {
			http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
			return
		}

		reqBodyBytes, err := io.ReadAll(req.Body)
		if checkErrorAsBadRequest(w, err) {
			return
		}

		reqData := new(struct {
			URL string `json:"url"`
		})
		err = json.Unmarshal(reqBodyBytes, reqData)
		if checkErrorAsBadRequest(w, err) {
			return
		}

		key, err := blogic.SaveNewURL(hs.stg, reqData.URL)
		if checkErrorAsInternalServerError(w, err) {
			return
		}

		respData := new(struct {
			Result string `json:"result"`
		})
		respData.Result = hs.cfg.BaseURL + "/" + key

		respBodyBytes, err := json.Marshal(respData)
		if checkErrorAsInternalServerError(w, err) {
			return
		}

		w.Header().Set("Content-Type", ApplicationJSONCharsetUTF8)
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write(respBodyBytes)
		if checkErrorAsInternalServerError(w, err) {
			return
		}
	}
}

func (hs *Handlers) makeGetURLHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		key := chi.URLParam(req, "key")

		if url, isOk := blogic.GetSavedURL(hs.stg, key); isOk {
			w.Header().Add("Location", url.String())
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "Resource Not Found", http.StatusNotFound)
		}
	}
}

func checkErrorAsBadRequest(w http.ResponseWriter, err error) bool {
	return checkError(w, err, http.StatusBadRequest)
}
func checkErrorAsInternalServerError(w http.ResponseWriter, err error) bool {
	return checkError(w, err, http.StatusInternalServerError)
}
func checkError(w http.ResponseWriter, err error, respHTTPCode int) bool {
	if err != nil {
		if respHTTPCode == 0 {
			respHTTPCode = http.StatusInternalServerError
		}
		http.Error(w, err.Error(), respHTTPCode)
		return true
	}
	return false
}
