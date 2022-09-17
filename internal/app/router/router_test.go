package router

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

func TestRouter(t *testing.T) {
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
			name:    "POST api correct #2",
			method:  http.MethodPost,
			request: "/api/shorten",
			body:    `{"url":"https://github.com/"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `{"result":"http://localhost:8080/3"}` + "\n",
				location:    "",
			},
		},
		{
			name:    "POST api correct #2",
			method:  http.MethodPost,
			request: "/api/shorten",
			body:    `{"url":"yandex.com/"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `{"result":"http://localhost:8080/4"}` + "\n",
				location:    "",
			},
		},
		{
			name:    "POST api incorrect #1",
			method:  http.MethodPost,
			request: "/api/shorten",
			body:    `{"url":"}`,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "error reading request\n",
				location:    "",
			},
		},
		{
			name:    "POST api incorrect #2",
			method:  http.MethodPost,
			request: "/api/shorten",
			body:    "https://github.com/",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "error reading request\n",
				location:    "",
			},
		},
		{
			name:    "POST correct link #1",
			method:  http.MethodPost,
			request: "/",
			body:    "https://github.com/",
			want: want{
				contentType: "text/plain;",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/5",
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
				body:        "http://localhost:8080/6",
				location:    "",
			},
		},
		{
			name:    "POST incorrect link #3",
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
				body:        "invalid id: 1sdG6\n",
				location:    "",
			},
		},
		{
			name:    "Invalid Method #1",
			method:  http.MethodGet,
			request: "/",
			body:    "",
			want: want{
				contentType: "",
				statusCode:  http.StatusMethodNotAllowed,
				body:        "",
				location:    "",
			},
		},
		{
			name:    "Invalid Method #2",
			method:  http.MethodPost,
			request: "/1993",
			body:    "",
			want: want{
				contentType: "",
				statusCode:  http.StatusMethodNotAllowed,
				body:        "",
				location:    "",
			},
		},
	}

	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080/"}
	mockRepo := storage.NewMockRepo()
	r := NewRouter(mockRepo, cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
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
