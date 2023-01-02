package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	type request struct {
		method  string
		target  string
		body    string
		headers map[string]string
	}
	type response struct {
		code         int
		hContentType string
		hLocation    string
		body         string
	}
	gzipStringFunc := func(str string) string {
		var buff bytes.Buffer
		gzw := gzip.NewWriter(&buff)
		_, err := gzw.Write([]byte(str))
		if err != nil {
			log.Fatal(err)
		}

		err = gzw.Close()
		if err != nil {
			log.Fatal(err)
		}
		return buff.String()
	}
	cfg := new(config.Config)
	config.LoadFromEnv(context.TODO(), cfg)
	tkits := []struct {
		key  string
		req  request
		resp response
	}{
		{
			key: "positive",
			req: request{
				method: http.MethodPost,
				target: "/",
				body:   "https://Kirill.Znamenskiy.me",
			},
			resp: response{
				code:         http.StatusCreated,
				hContentType: "text/plain;charset=UTF-8",
				body:         `^` + cfg.BaseURL + `/[-\w]+$`,
			},
		},
		{
			key: "positive",
			req: request{
				method:  http.MethodPost,
				target:  "/",
				headers: map[string]string{"Accept-Encoding": "gzip"},
				body:    "https://Kirill.Znamenskiy.me",
			},
			resp: response{
				code:         http.StatusCreated,
				hContentType: "text/plain;charset=UTF-8",
			},
		},
		{
			key: "positive",
			req: request{
				method: http.MethodPost,
				target: "/",
				headers: map[string]string{
					"Content-Encoding": "gzip",
					"Accept-Encoding":  "gzip",
					"Content-Type":     "application/x-gzip",
				},
				body: gzipStringFunc("https://Kirill.Znamenskiy.me"),
			},
			resp: response{
				code:         http.StatusCreated,
				hContentType: "text/plain;charset=UTF-8",
			},
		},
		{
			key: "negative",
			req: request{
				method: http.MethodPost,
				target: "/api/shorten",
				body:   `{"URL": "https://Kirill.Znamenskiy.pw"}`,
			},
			resp: response{
				code: http.StatusBadRequest,
			},
		},
		{
			key: "positive",
			req: request{
				method:  http.MethodPost,
				target:  "/api/shorten",
				body:    `{"URL": "https://Kirill.Znamenskiy.pw"}`,
				headers: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
			},
			resp: response{
				code:         http.StatusCreated,
				hContentType: "application/json;charset=UTF-8",
				body:         `^\{\"result\"\:\"` + cfg.BaseURL + `/[-\w]+\"\}$`,
			},
		},
		{
			key: "positive",
			req: request{
				method: http.MethodGet,
				target: "/positive-test-2",
			},
			resp: response{
				code:      http.StatusTemporaryRedirect,
				hLocation: "https://Kirill.Znamenskiy.pw",
			},
		},
		{
			key: "negative",
			req: request{
				method: http.MethodHead,
				target: "/",
			},
			resp: response{
				code: http.StatusBadRequest,
			},
		},
		{
			key: "negative",
			req: request{
				method: http.MethodGet,
				target: "/",
			},
			resp: response{
				code: http.StatusBadRequest,
			},
		},
		{
			key: "negative",
			req: request{
				method: http.MethodGet,
				target: "/adgg",
			},
			resp: response{
				code: http.StatusNotFound,
			},
		},
		{
			key: "negative",
			req: request{
				method: http.MethodPost,
				target: "/",
				body:   ":ht3240dfkk",
			},
			resp: response{
				code: http.StatusInternalServerError,
			},
		},
	}
	stg := storage.NewInMemoryStorage()
	u, err := url.Parse("https://Kirill.Znamenskiy.pw")
	if err != nil {
		t.Fatal(err)
	}
	err = stg.Put("positive-test-2", u)
	if err != nil {
		t.Fatal(err)
	}
	cfg.SetStorage(stg)
	for tind, tkit := range tkits {
		t.Run(fmt.Sprintf("Test %d %s", tind+1, tkit.key), func(t *testing.T) {
			var tstReqBody io.Reader
			if tkit.req.body != "" {
				tstReqBody = strings.NewReader(tkit.req.body)
			}

			// создаём новый Recorder
			w := httptest.NewRecorder()

			req := httptest.NewRequest(tkit.req.method, tkit.req.target, tstReqBody)
			//if tkit.req.headers == nil {
			//	tkit.req.headers = map[string]string{"Content-Type": "text/html;charset=UTF-8"}
			//}
			for hName, hValue := range tkit.req.headers {
				req.Header.Set(hName, hValue)
			}

			// определяем хендлер
			h := MakeMainHandler(cfg)
			// запускаем сервер
			h.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			// проверяем код ответа
			if resp.StatusCode != tkit.resp.code {
				t.Errorf("Expected status code %d, got %d", tkit.resp.code, w.Code)
			}

			if tkit.resp.hContentType != "" && resp.Header.Get("Content-Type") != tkit.resp.hContentType {
				t.Errorf("Expected Content-Type %s, got %s", tkit.resp.hContentType, resp.Header.Get("Content-Type"))
			}
			if tkit.resp.hLocation != "" && resp.Header.Get("Location") != tkit.resp.hLocation {
				t.Errorf("Expected Location %s, got %s", tkit.resp.hLocation, resp.Header.Get("Location"))
			}

			if tkit.resp.body != "" {
				respBodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				respBodyString := string(respBodyBytes)
				if rgxp, err := regexp.Compile(tkit.resp.body); err == nil {
					if !rgxp.Match(respBodyBytes) {
						t.Errorf("Expected Body match pattern %s, got %s", tkit.resp.body, respBodyBytes)
					}
				} else {
					if respBodyString != tkit.resp.body {
						t.Errorf("Expected Body %s, got %s", tkit.resp.body, respBodyBytes)
					}
				}
			}
		})
	}
}