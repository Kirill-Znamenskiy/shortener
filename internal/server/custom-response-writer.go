package server

import (
	"net/http"
)

type CustomResponseWriter struct {
	Body       []byte
	StatusCode int
	header     http.Header
}

func NewCustomResponseWriter() *CustomResponseWriter {
	return &CustomResponseWriter{
		header: http.Header{},
	}
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.header
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.Body = b
	return len(b), nil
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func (w *CustomResponseWriter) WriteToNext(nw http.ResponseWriter) (int, error) {
	for hKey := range nw.Header() {
		nw.Header().Del(hKey)
	}
	for hKey, hValues := range w.Header() {
		for _, hValue := range hValues {
			nw.Header().Add(hKey, hValue)
		}
	}
	nw.WriteHeader(w.StatusCode)

	return nw.Write(w.Body)
}
