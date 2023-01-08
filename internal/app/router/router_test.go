package router

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type want struct {
	contentType string
	statusCode  int
	body        string
	location    string
}

type test struct {
	name        string
	compression string
	method      string
	request     string
	body        string
	want        want
}

func runRouterTest(t *testing.T, tests []test, db bool) {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	var mockRepo storage.Repository
	var prefix string
	var err error
	switch {
	case db:
		prefix = "Postgres repo: "
		dbURL := os.Getenv("DATABASE_DSN")
		if dbURL == "" {
			return
		}
		mockRepo, err = storage.NewPostgresMockRepo(dbURL)
		require.NoError(t, err)
		defer mockRepo.DeleteRepo(context.Background())
	default:
		mockRepo = storage.NewMockRepo()
		prefix = "InMem repo: "
	}
	r := NewRouter(mockRepo, cfg)
	for _, tt := range tests {
		t.Run(prefix+tt.name, func(t *testing.T) {
			var reader io.ReadCloser
			var err error
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			request.Header.Set("Accept-Encoding", tt.compression)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			switch tt.compression {
			case "gzip":
				reader, err = gzip.NewReader(result.Body)
				require.NoError(t, err)
			default:
				reader = result.Body
			}
			bodyResult, err := io.ReadAll(reader)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(bodyResult))
		})
	}
}

func TestPostApiCorrect(t *testing.T) {
	tests := []test{
		{name: "POST api correct #1",
			method:  http.MethodPost,
			request: "/api/shorten",
			body:    `{"url":"https://github.com/"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `{"result":"http://localhost:8080/vRveliyDLz8"}` + "\n",
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
				body:        `{"result":"http://localhost:8080/gjsBFlccqF6"}` + "\n",
				location:    "",
			},
		},
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestPostApiIncorrect(t *testing.T) {
	tests := []test{
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
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestPostCorrect(t *testing.T) {
	tests := []test{
		{
			name:    "POST correct link #1",
			method:  http.MethodPost,
			request: "/",
			body:    "https://github.com/",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/vRveliyDLz8",
				location:    "",
			},
		},
		{
			name:    "POST correct link #2",
			method:  http.MethodPost,
			request: "/",
			body:    "https://yandex.ru/1234",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/cXRXuMGP3pD",
				location:    "",
			},
		},
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestPostIncorrect(t *testing.T) {
	tests := []test{
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
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestGet(t *testing.T) {
	tests := []test{
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
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestMethod(t *testing.T) {
	tests := []test{
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
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestApiPostCompressed(t *testing.T) {
	tests := []test{
		{
			name:        "POST api correct with compression",
			compression: "gzip",
			method:      http.MethodPost,
			request:     "/api/shorten",
			body:        `{"url":"https://github.com/"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `{"result":"http://localhost:8080/vRveliyDLz8"}` + "\n",
				location:    "",
			},
		},
		{
			name:        "POST api incorrect with compression",
			compression: "gzip",
			method:      http.MethodPost,
			request:     "/api/shorten",
			body:        "https://github.com/",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "error reading request\n",
				location:    "",
			},
		},
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestBatchCorrect(t *testing.T) {
	tests := []test{
		{
			name:        "Correct batch POST api",
			compression: "gzip",
			method:      http.MethodPost,
			request:     "/api/shorten/batch",
			body:        `[{"correlation_id":"TestCorrelationID1","original_url":"https://github.com/"},{"correlation_id":"TestCorrelationID2","original_url":"https://yandex.com/"}]`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `[{"correlation_id":"TestCorrelationID1","short_url":"http://localhost:8080/vRveliyDLz8"},{"correlation_id":"TestCorrelationID2","short_url":"http://localhost:8080/BlbEuA4l5GJ"}]` + "\n",
				location:    "",
			},
		},
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestBatchIncorrect(t *testing.T) {
	tests := []test{
		{
			name:        "Incorrect batch POST api",
			compression: "gzip",
			method:      http.MethodPost,
			request:     "/api/shorten/batch",
			body:        "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "error reading request\n",
				location:    "",
			},
		},
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestPostDuplicate(t *testing.T) {
	tests := []test{
		{
			name:    "POST correct link #1",
			method:  http.MethodPost,
			request: "/",
			body:    "https://github.com/",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/vRveliyDLz8",
				location:    "",
			},
		},
		{
			name:    "POST correct link #2",
			method:  http.MethodPost,
			request: "/",
			body:    "https://github.com/",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusConflict,
				body:        "http://localhost:8080/vRveliyDLz8",
				location:    "",
			},
		},
		{
			name:    "POST correct link #3",
			method:  http.MethodPost,
			request: "/",
			body:    "https://github.com/",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusConflict,
				body:        "http://localhost:8080/vRveliyDLz8",
				location:    "",
			},
		},
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestApiDuplicate(t *testing.T) {
	tests := []test{
		{
			name:        "POST DB add duplicate #1",
			compression: "gzip",
			method:      http.MethodPost,
			request:     "/api/shorten",
			body:        `{"url":"https://github.com/"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `{"result":"http://localhost:8080/vRveliyDLz8"}` + "\n",
				location:    "",
			},
		},
		{
			name:        "POST DB add duplicate #2",
			compression: "gzip",
			method:      http.MethodPost,
			request:     "/api/shorten",
			body:        `{"url":"https://github.com/"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusConflict,
				body:        `{"result":"http://localhost:8080/vRveliyDLz8"}` + "\n",
				location:    "",
			},
		},
		{
			name:        "POST DB add duplicate #3",
			compression: "gzip",
			method:      http.MethodPost,
			request:     "/api/shorten",
			body:        `{"url":"https://github.com/"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusConflict,
				body:        `{"result":"http://localhost:8080/vRveliyDLz8"}` + "\n",
				location:    "",
			},
		},
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}

func TestApiBatchDuplicate(t *testing.T) {
	tests := []test{
		{
			name:        "Correct batch duplicate POST api.",
			compression: "gzip",
			method:      http.MethodPost,
			request:     "/api/shorten/batch",
			body:        `[{"correlation_id":"TestCorrelationID1","original_url":"https://github.com/"},{"correlation_id":"TestCorrelationID2","original_url":"https://github.com/"}]`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusConflict,
				body:        `[{"correlation_id":"TestCorrelationID1","short_url":"http://localhost:8080/vRveliyDLz8"},{"correlation_id":"TestCorrelationID2","short_url":"http://localhost:8080/vRveliyDLz8"}]` + "\n",
				location:    "",
			},
		},
	}
	runRouterTest(t, tests, true)
	runRouterTest(t, tests, false)
}
