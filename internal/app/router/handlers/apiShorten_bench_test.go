package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func BenchmarkShortenAPI(b *testing.B) {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	dbURL := os.Getenv("DATABASE_DSN")
	repo, err := storage.NewPostgresMockRepo(dbURL)
	require.NoError(b, err)
	defer repo.DeleteRepo(context.Background())

	handler := APIShorten(repo, cfg)
	b.ResetTimer()
	b.Run("APIShorten", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			body, _ := json.Marshal(model.URL{LongURL: fmt.Sprint(i) + ".ru"})
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			b.StartTimer()
			handler.ServeHTTP(w, request)
			_ = w.Result().Body.Close()
		}
	})

}
