package handlers

import (
	"context"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func BenchmarkAPIUserExpandAPI(b *testing.B) {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	dbURL := os.Getenv("DATABASE_DSN")
	repo, err := storage.NewPostgresMockRepo(dbURL)
	require.NoError(b, err)
	defer repo.DeleteRepo(context.Background())

	handler := APIUserExpand(repo, cfg)
	urls := make(map[string]*model.URL, 10000)
	for i := 0; i < 10000; i++ {
		urls[fmt.Sprint(i)] = &model.URL{UserID: "testUser",
			ShortURL: encoders.ToRBase62(fmt.Sprint(i)),
			LongURL:  fmt.Sprint(i) + ".com"}
	}
	repo.AddBatch(context.Background(), urls)

	b.ResetTimer()
	b.Run("ExpandUserAPI", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			w := httptest.NewRecorder()
			http.SetCookie(w, &http.Cookie{Name: "user_id", Value: "testUser"})
			b.StartTimer()
			handler.ServeHTTP(w, request)
			_ = w.Result().Body.Close()
		}
	})
}
