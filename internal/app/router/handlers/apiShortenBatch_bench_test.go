package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func BenchmarkAPIShortenBatch(b *testing.B) {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	dbURL := os.Getenv("DATABASE_DSN")
	repo, err := storage.NewPostgresMockRepo(dbURL)
	defer repo.DeleteRepo(context.Background())
	require.NoError(b, err)

	handler := APIShortenBatch(repo, cfg)
	ShortenItems := make([]model.BatchReqItem, 0, 1000)
	for i := 0; i < cap(ShortenItems); i++ {
		ShortenItems = append(ShortenItems, model.BatchReqItem{CorID: fmt.Sprint(i), OrigURL: uuid.NewString() + ".ru"})
	}
	b.ResetTimer()
	b.Run("ShortenAPIBatch", func(b *testing.B) {
		b.StopTimer()
		body, err := json.Marshal(ShortenItems)
		require.NoError(b, err)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		b.StartTimer()
		handler.ServeHTTP(w, request)
		_ = w.Result().Body.Close()
	})
}
