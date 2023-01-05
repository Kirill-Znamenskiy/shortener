package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/types"
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
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

	shortener := blogic.NewShortener(cfg.BaseURL, cfg.GetStorage())
	hs := &handlers{cfg: cfg, shortener: shortener}

	r := chi.NewRouter()

	r.Use(CleanURLPathMiddleware())
	r.Use(AllowContentCharsetMiddleware("", UTF8))

	r.Use(DecompressMiddleware())
	r.Use(middleware.Compress(5, TextHTML, TextPlain, ApplicationJSON))

	r.Use(middleware.AllowContentType("", TextHTML, TextPlain, ApplicationJSON))

	r.Use(AuthUserMiddleware(cfg, shortener.GenerateNewUser))

	r.Post("/", hs.makeWrapperForJSONHandlerFunc(hs.makeSaveNewURLHandlerFunc()))
	r.Get("/{key:[-\\w]+}", hs.makeGetURLHandlerFunc())

	r.Post("/api/shorten", hs.makeSaveNewURLHandlerFunc())
	r.Get("/api/user/urls", hs.makeGetUserURLsHandlerFunc())

	r.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("UNEXPECTED: %s %q\n", req.Method, req.URL.Path)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	return r
}

type handlers struct {
	cfg       *config.Config
	shortener *blogic.Shortener
}

func (hs *handlers) makeWrapperForJSONHandlerFunc(nextJSONHandler http.Handler) http.HandlerFunc {
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
			reqBodyCont = []byte(`{"url":"` + string(reqBodyCont) + `"}`)
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

func (hs *handlers) makeSaveNewURLHandlerFunc() http.HandlerFunc {
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

		userUUID, err := extractUser(req)
		if checkErrorAsInternalServerError(w, err) {
			return
		}

		record, err := hs.shortener.SaveNewURL(userUUID, reqData.URL)
		if checkErrorAsInternalServerError(w, err) {
			return
		}

		respData := new(struct {
			Result string `json:"result"`
		})
		respData.Result = hs.shortener.BuildShortURL(record)

		finishHandler(w, respData, http.StatusCreated)
	}
}

func (hs *handlers) makeGetURLHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		key := chi.URLParam(req, "key")

		if url, isOk := hs.shortener.GetSavedURL(key); isOk {
			w.Header().Add("Location", url.String())
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "Resource Not Found", http.StatusNotFound)
		}
	}
}

func (hs *handlers) makeGetUserURLsHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, err := extractUser(req)
		if checkErrorAsInternalServerError(w, err) {
			return
		}
		allUserRecords := hs.shortener.GetAllUserRecords(user)

		if len(allUserRecords) == 0 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			respData := make([]struct {
				ShortURL    string `json:"short_url"`
				OriginalURL string `json:"original_url"`
			}, len(allUserRecords))
			ind := 0
			for _, record := range allUserRecords {
				respData[ind].ShortURL = hs.shortener.BuildShortURL(record)
				respData[ind].OriginalURL = record.OriginalURL.String()
			}

			finishHandler(w, respData, http.StatusOK)
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

type myString string

const userContextValueKey myString = "user"

func extractUser(req *http.Request) (ret types.User, err error) {
	ret, isOk := req.Context().Value(userContextValueKey).(types.User)
	if !isOk {
		return nil, errors.New("error at extracting user")
	}
	return
}

func finishHandler(w http.ResponseWriter, respData any, respStatusCode int) {

	respBodyBytes, err := json.Marshal(respData)
	if checkErrorAsInternalServerError(w, err) {
		return
	}

	w.Header().Set("Content-Type", ApplicationJSONCharsetUTF8)

	if respStatusCode == 0 {
		respStatusCode = http.StatusOK
	}
	w.WriteHeader(respStatusCode)

	_, err = w.Write(respBodyBytes)
	if checkErrorAsInternalServerError(w, err) {
		return
	}
}
