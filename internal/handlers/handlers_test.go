package handlers

import (
	"bytes"
	"compress/gzip"
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
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
	cfg := config.LoadFromEnv()
	tests := []struct {
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
				code: http.StatusBadRequest,
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
				code: http.StatusBadRequest,
			},
		},
	}
	stg := storage.NewInMemoryStorage()
	u, err := url.Parse("https://Kirill.Znamenskiy.pw")
	if err != nil {
		t.Fatal(err)
	}
	_, err = stg.Put("positive-test-2", u)
	if err != nil {
		t.Fatal(err)
	}
	for tstInd, tst := range tests {
		t.Run("Test "+strconv.Itoa(tstInd+1)+" "+tst.key, func(t *testing.T) {
			var tstReqBody io.Reader
			if tst.req.body != "" {
				tstReqBody = strings.NewReader(tst.req.body)
			}

			// ?????????????? ?????????? Recorder
			w := httptest.NewRecorder()

			req := httptest.NewRequest(tst.req.method, tst.req.target, tstReqBody)
			//if tst.req.headers == nil {
			//	tst.req.headers = map[string]string{"Content-Type": "text/html;charset=UTF-8"}
			//}
			for hName, hValue := range tst.req.headers {
				req.Header.Set(hName, hValue)
			}

			// ???????????????????? ??????????????
			h := MakeMainHandler(stg, cfg)
			// ?????????????????? ????????????
			h.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			// ?????????????????? ?????? ????????????
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
						t.Errorf("Expected Body match pattern %s, got %s", tst.resp.body, respBodyBytes)
					}
				} else {
					if respBodyString != tst.resp.body {
						t.Errorf("Expected Body %s, got %s", tst.resp.body, respBodyBytes)
					}
				}
			}
		})
	}
}
