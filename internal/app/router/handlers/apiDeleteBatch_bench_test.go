package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/router/loader"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func BenchmarkAPIDeleteBatchAPI(b *testing.B) {
	dbURL := os.Getenv("DATABASE_DSN")
	repo, err := storage.NewPostgresMockRepo(dbURL)
	require.NoError(b, err)
	defer repo.DeleteRepo(context.Background())
	handler := APIDeleteBatch(loader.NewDeleteLoader(repo))
	urls := make(map[string]*models.URL, 10000)
	deleteURLs := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		urls[fmt.Sprint(i)] = &models.URL{UserID: "user1",
			ShortURL: encoders.ToRBase62(fmt.Sprint(i)),
			LongURL:  fmt.Sprint(i) + ".com"}
	}
	_, err = repo.AddBatch(context.Background(), urls)
	require.NoError(b, err)
	for i := 0; i < 10000; i++ {
		deleteURLs = append(deleteURLs, encoders.ToRBase62(fmt.Sprint(i)))
	}
	b.ResetTimer()
	b.Run("APIDeleteBatch", func(b *testing.B) {
		b.StopTimer()
		body, err := json.Marshal(deleteURLs)
		require.NoError(b, err)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(body)))
		request.Header = map[string][]string{"Cookie": {"user_id=user1", "signature=60e8d0babc58e796ac223a64b5e68b998de7d3b203bc8a859bc0ec15ee66f5f9"}}
		w := httptest.NewRecorder()
		b.StartTimer()
		handler.ServeHTTP(w, request)
		_ = w.Result().Body.Close()
	})
}
