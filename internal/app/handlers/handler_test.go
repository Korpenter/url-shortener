package handlers

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenerHandler_ServeHTTP(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
		location    string
	}
	tests := []struct {
		name    string
		method  string
		request string
		body    string
		want    want
	}{
		{
			name:    "POST correct link #1",
			method:  http.MethodPost,
			request: "/",
			body:    "https://github.com/",
			want: want{
				contentType: "text/plain;",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/1",
				location:    "",
			},
		},
		{
			name:    "POST correct link #2",
			method:  http.MethodPost,
			request: "/",
			body:    "https://yandex.ru/",
			want: want{
				contentType: "text/plain;",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/2",
				location:    "",
			},
		},
		{
			name:    "POST incorrect link ",
			method:  http.MethodPost,
			request: "/",
			body:    "https://",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "invalid url\n",
				location:    "",
			},
		},
		{
			name:    "GET present id ",
			method:  http.MethodGet,
			request: "/2",
			body:    "",
			want: want{
				contentType: "",
				statusCode:  http.StatusTemporaryRedirect,
				body:        "",
				location:    "https://yandex.ru/",
			},
		},
		{
			name:    "GET invalid id ",
			method:  http.MethodGet,
			request: "/1sdG6",
			body:    "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				body:        "invalid id\n",
				location:    "",
			},
		},
		{
			name:    "GET empty id ",
			method:  http.MethodGet,
			request: "/",
			body:    "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "no id\n",
				location:    "",
			},
		},
	}
	cfg := config.New()
	urls := storage.NewInMemRepo()
	handler := NewShortenerHandler(urls, cfg)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))

			bodyResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(bodyResult))
		})
	}
}
