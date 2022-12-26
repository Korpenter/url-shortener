package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/stretchr/testify/require"
)

func BenchmarkExpand(b *testing.B) {
	dbURL := os.Getenv("DATABASE_DSN")
	var repo storage.Repository
	var err error
	if dbURL != "" {
		repo, err = storage.NewPostgresMockRepo(dbURL)
		require.NoError(b, err)
	} else {
		repo = storage.NewMockRepo()
	}
	urls := make(map[string]*models.URL, 10000)
	for i := 0; i < 10000; i++ {
		urls[fmt.Sprint(i)] = &models.URL{UserID: "user1",
			ShortURL: encoders.ToRBase62(fmt.Sprint(i)),
			LongURL:  fmt.Sprint(i) + ".ru"}
	}
	repo.AddBatch(context.Background(), urls)
	handler := Expand(repo)
	b.ResetTimer()
	b.Run("Expand", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.ReportAllocs()
			b.StopTimer()
			request := httptest.NewRequest(http.MethodGet, "/"+encoders.ToRBase62(fmt.Sprint(i)), nil)
			request.Header = map[string][]string{"Cookie": {"user_id=user1", "signature=60e8d0babc58e796ac223a64b5e68b998de7d3b203bc8a859bc0ec15ee66f5f9"}}
			w := httptest.NewRecorder()
			b.StartTimer()
			handler.ServeHTTP(w, request)
			_ = w.Result().Body.Close()
		}
	})
}
