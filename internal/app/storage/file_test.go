package storage

import (
	"testing"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileRepo(t *testing.T) {
	tests := []struct {
		name        string
		urls        []models.URL
		wantContain []models.URL
	}{
		{
			name: "Load preexisting file, No addition",
			wantContain: []models.URL{
				{ShortURL: "1", LongURL: "yandex.ru", UserID: "Helloworld"},
				{ShortURL: "2", LongURL: "hero.ru", UserID: "Helloworld"},
			},
		},
		{
			name: " Add links and Load preexisting file.",
			urls: []models.URL{
				{ShortURL: "3", LongURL: "hell.ru", UserID: "Anotherone"},
				{ShortURL: "4", LongURL: "nvidia.ru", UserID: "Anotherone"},
			},
			wantContain: []models.URL{
				{ShortURL: "3", LongURL: "hell.ru", UserID: "Anotherone"},
				{ShortURL: "4", LongURL: "nvidia.ru", UserID: "Anotherone"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fRepo, err := newMockFileRepo()
			require.NoError(t, err)
			for _, v := range tt.urls {
				_, err = fRepo.add(v.LongURL, v.ShortURL, v.UserID)
				require.NoError(t, err)
			}
			err = fRepo.load()
			require.NoError(t, err)
			var urls []models.URL
			for _, value := range fRepo.cacheByShort {
				urls = append(urls, *value)
			}
			for _, v := range tt.wantContain {
				assert.Contains(t, urls, v)

			}
			err = fRepo.delete()
			require.NoError(t, err)
		})
	}
}
