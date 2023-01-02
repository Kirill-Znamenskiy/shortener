package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"path"
	"strconv"
)

func DecompressMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				// skip check for empty content body
				next.ServeHTTP(w, r)
				return
			}

			switch ce := r.Header.Get("Content-Encoding"); ce {
			case "":
			case "gzip":
				reader, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				reqBodyBytes, err := io.ReadAll(reader)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = r.Body.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				r.Body = io.NopCloser(bytes.NewReader(reqBodyBytes))
				r.ContentLength = int64(len(reqBodyBytes))
				r.Header.Set("Content-Length", strconv.FormatInt(r.ContentLength, 10))
				r.Header.Del("Content-Encoding")
				r.Header.Del("Content-Type")
			default:
				http.Error(w, fmt.Sprintf("Unknown request Content-Encoding: %q", ce), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AllowContentCharsetMiddleware generates a handler that writes a 415 Unsupported Media Type response if none of
// the charsets match.
// An empty charset will allow requests with no Content-Type header or no specified charset.
// Skip check for empty content body.
// Basically this is just a wrapper to chi middleware ContentCharset, with check bout empty content body.
func AllowContentCharsetMiddleware(charsets ...string) func(next http.Handler) http.Handler {
	chiMiddleware := middleware.ContentCharset(charsets...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				// skip check for empty content body
				next.ServeHTTP(w, r)
				return
			}

			chiMiddleware(next).ServeHTTP(w, r)
		})
	}
}

func CleanURLPathMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := path.Clean(r.URL.Path)
			if res == "" || res == "." {
				res = "/"
			}
			r.URL.Path = res

			next.ServeHTTP(w, r)
		})
	}
}
