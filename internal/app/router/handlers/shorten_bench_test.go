package handlers

import (
	"context"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func BenchmarkShorten(b *testing.B) {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	dbURL := os.Getenv("DATABASE_DSN")
	repo, err := storage.NewPostgresMockRepo(dbURL)
	require.NoError(b, err)
	defer repo.DeleteRepo(context.Background())

	handler := Shorten(repo, cfg)
	b.ResetTimer()
	b.Run("Shorten", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprint(i)+".ru"))
			w := httptest.NewRecorder()
			b.StartTimer()
			handler.ServeHTTP(w, request)
			_ = w.Result().Body.Close()
		}
	})
}
