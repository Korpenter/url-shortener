package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/stretchr/testify/require"
)

func BenchmarkAPIUserExpandAPI(b *testing.B) {
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
	shortener := service.NewShortenerImpl(repo, cfg)
	handler := APIUserExpand(shortener)
	urls := make([]*models.URL, 10000)
	for i := 0; i < 10000; i++ {
		urls[i] = &models.URL{UserID: "user1",
			ShortURL: encoders.ToRBase62(fmt.Sprint(i)),
			LongURL:  fmt.Sprint(i) + ".com"}
	}
	repo.AddBatch(context.Background(), urls)

	b.ResetTimer()
	b.Run("ExpandUserAPI", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			request.Header = map[string][]string{"Cookie": {"user_id=user1", "signature=60e8d0babc58e796ac223a64b5e68b998de7d3b203bc8a859bc0ec15ee66f5f9"}}
			w := httptest.NewRecorder()
			b.StartTimer()
			handler.ServeHTTP(w, request)
			_ = w.Result().Body.Close()
		}
	})
}
