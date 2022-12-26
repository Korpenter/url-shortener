package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func BenchmarkAPIShortenBatch(b *testing.B) {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	dbURL := os.Getenv("DATABASE_DSN")
	var repo storage.Repository
	var err error
	if dbURL != "" {
		repo, err = storage.NewPostgresMockRepo(dbURL)
		require.NoError(b, err)
	} else {
		repo = storage.NewMockRepo()
	}

	handler := APIShortenBatch(repo, cfg)
	ShortenItems := make([]models.BatchReqItem, 0, 10000)
	for i := 0; i < cap(ShortenItems); i++ {
		ShortenItems = append(ShortenItems, models.BatchReqItem{CorID: fmt.Sprint(i), OrigURL: uuid.NewString() + ".ru"})
	}
	b.ResetTimer()
	b.Run("ShortenAPIBatch", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			body, err := json.Marshal(ShortenItems)
			require.NoError(b, err)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(body)))
			request.Header = map[string][]string{"Cookie": {"user_id=user1", "signature=60e8d0babc58e796ac223a64b5e68b998de7d3b203bc8a859bc0ec15ee66f5f9"}}
			w := httptest.NewRecorder()
			b.StartTimer()
			handler.ServeHTTP(w, request)
			_ = w.Result().Body.Close()
		}
	})
}
