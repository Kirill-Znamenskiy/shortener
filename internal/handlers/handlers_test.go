package handlers

import (
	"github.com/Kirill-Znamenskiy/shortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	type request struct {
		method string
		target string
		body   string
	}
	type response struct {
		code         int
		hContentType string
		hLocation    string
		body         string
	}
	tests := []struct {
		name string
		req  request
		resp response
	}{
		{
			name: "positive test #1",
			req: request{
				method: http.MethodPost,
				target: "/",
				body:   "https://Kirill.Znamenskiy.pw",
			},
			resp: response{
				code:         201,
				hContentType: "text/plain; charset=UTF-8",
				body:         `^http://localhost:8080/[-\w]+$`,
			},
		},
		{
			name: "positive test #2",
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
			name: "negative test #1",
			req: request{
				method: http.MethodHead,
				target: "/",
			},
			resp: response{
				code: 400,
			},
		},
		{
			name: "negative test #2",
			req: request{
				method: http.MethodGet,
				target: "/",
			},
			resp: response{
				code: 400,
			},
		},
		{
			name: "negative test #3",
			req: request{
				method: http.MethodGet,
				target: "/adgg",
			},
			resp: response{
				code: 400,
			},
		},
		{
			name: "negative test #4",
			req: request{
				method: http.MethodPost,
				target: "/",
				body:   ":ht3240dfkk",
			},
			resp: response{
				code: 400,
			},
		},
	}
	stg := storage.NewInMemoryStorage()
	u, err := url.Parse("https://Kirill.Znamenskiy.pw")
	if err != nil {
		t.Fatal(err)
	}
	stg.Put("positive-test-2", u)
	hs := Handlers{Stg: stg}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			var tstReqBody io.Reader
			if tst.req.body != "" {
				tstReqBody = strings.NewReader(tst.req.body)
			}

			// создаём новый Recorder
			w := httptest.NewRecorder()

			req := httptest.NewRequest(tst.req.method, tst.req.target, tstReqBody)

			// определяем хендлер
			h := hs.MakeMainHandler()
			// запускаем сервер
			h.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			// проверяем код ответа
			if resp.StatusCode != tst.resp.code {
				t.Errorf("Expected status code %d, got %d", tst.resp.code, w.Code)
			}

			if tst.resp.hContentType != "" && resp.Header.Get("Content-Type") != tst.resp.hContentType {
				t.Errorf("Expected Content-Type %s, got %s", tst.resp.hContentType, resp.Header.Get("Content-Type"))
			}
			if tst.resp.hLocation != "" && resp.Header.Get("Location") != tst.resp.hLocation {
				t.Errorf("Expected Location %s, got %s", tst.resp.hLocation, resp.Header.Get("Location"))
			}

			if tst.resp.body != "" {
				respBodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				respBodyString := string(respBodyBytes)
				if rgxp, err := regexp.Compile(tst.resp.body); err == nil {
					if !rgxp.Match(respBodyBytes) {
						t.Errorf("Expected body match pattern %s, got %s", tst.resp.body, respBodyBytes)
					}
				} else {
					if respBodyString != tst.resp.body {
						t.Errorf("Expected body %s, got %s", tst.resp.body, respBodyBytes)
					}
				}
			}
		})
	}
}
