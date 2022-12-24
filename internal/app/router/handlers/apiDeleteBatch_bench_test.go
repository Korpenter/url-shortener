package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/model"
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
	urls := make(map[string]*model.URL, 10000)
	deleteURLs := make([]*model.URL, 10000)
	for i := 0; i < 10000; i++ {
		urls[fmt.Sprint(i)] = &model.URL{UserID: "testUser",
			ShortURL: encoders.ToRBase62(fmt.Sprint(i)),
			LongURL:  fmt.Sprint(i) + ".com"}
	}
	_, err = repo.AddBatch(context.Background(), urls)
	require.NoError(b, err)
	for i := 0; i < 10000; i++ {
		deleteURLs = append(deleteURLs, &model.URL{
			ShortURL: encoders.ToRBase62(fmt.Sprint(i)),
		})
	}

	b.ResetTimer()
	b.Run("APIDeleteBatch", func(b *testing.B) {
		b.StopTimer()
		body, err := json.Marshal(deleteURLs)
		require.NoError(b, err)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		b.StartTimer()
		handler.ServeHTTP(w, request)
		_ = w.Result().Body.Close()
	})
}
