package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mldlr/url-shortener/internal/app/models"
)

func TestInMemRepo_Add(t *testing.T) {
	tests := []struct {
		name        string
		longURL     string
		userID      string
		wantShort   string
		wantContain models.URL
	}{
		{
			name:      "Successfully added link",
			longURL:   "https://github.com/",
			userID:    "KS097f1lS&F",
			wantShort: "vRveliyDLz8",
			wantContain: models.URL{
				ShortURL: "vRveliyDLz8",
				LongURL:  "https://github.com/",
				UserID:   "KS097f1lS&F",
			},
		},
		{
			name:      "Successfully added link",
			longURL:   "https://github.com/1234",
			userID:    "KS097f1lS&F",
			wantShort: "7xhhSdE3RuA",
			wantContain: models.URL{
				ShortURL: "7xhhSdE3RuA",
				LongURL:  "https://github.com/1234",
				UserID:   "KS097f1lS&F",
			},
		},
	}
	mock := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := mock.NewID(tt.longURL)
			url := models.URL{ShortURL: id, LongURL: tt.longURL, UserID: tt.userID}
			_, _ = mock.Add(context.Background(), &url)
			assert.Equal(t, tt.wantShort, url.ShortURL)
			var urls []models.URL
			for _, value := range mock.urlsByShort {
				urls = append(urls, *value)
			}
			assert.Contains(t, urls, tt.wantContain)
			for _, value := range mock.urlsByUser[tt.userID] {
				urls = append(urls, *value)
			}
			assert.Contains(t, urls, tt.wantContain)
		})
	}
}

func TestInMemRepo_GetByShort(t *testing.T) {
	tests := []struct {
		name    string
		short   string
		want    *models.URL
		wantErr bool
	}{
		{
			name:    "Id in repo",
			short:   "aQqomlSbUsE",
			want:    &models.URL{ShortURL: "aQqomlSbUsE", LongURL: "https://yandex.ru/", UserID: "KS097f1lS&F", Deleted: false},
			wantErr: false,
		},
		{
			name:    "Id not in repo",
			short:   "3",
			want:    nil,
			wantErr: true,
		},
	}
	mock := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mockRepo{
				urlsByShort: mock.urlsByShort,
				urlsByUser:  mock.urlsByUser,
			}
			got, err := r.Get(context.Background(), tt.short)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestInMemRepo_GetByUser(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		want    []*models.URL
		wantErr bool
	}{
		{
			name:   "User has 2 urls",
			userID: "KS097f1lS&F",
			want: []*models.URL{
				{
					ShortURL: "3S93m80EGmF",
					LongURL:  "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders",
					UserID:   "KS097f1lS&F",
				},

				{
					ShortURL: "aQqomlSbUsE",
					LongURL:  "https://yandex.ru/",
					UserID:   "KS097f1lS&F",
				},
			},
			wantErr: false,
		},
		{
			name:    "User has no urls",
			userID:  "SDADAD&FS()AS",
			want:    nil,
			wantErr: false,
		},
	}
	mock := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockRepo{
				urlsByShort: mock.urlsByShort,
				urlsByUser:  mock.urlsByUser,
			}
			got, err := r.GetByUser(context.Background(), tt.userID)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
				return
			}
			assert.Error(t, err)
		})
	}
}
