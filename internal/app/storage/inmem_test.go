package storage

import (
	"testing"

	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemRepo_Add(t *testing.T) {
	tests := []struct {
		name        string
		longURL     string
		userID      string
		wantShort   string
		wantContain model.URL
	}{
		{
			name:      "Successfully added link",
			longURL:   "https://github.com/",
			userID:    "KS097f1lS&F",
			wantShort: "3",
			wantContain: model.URL{
				ShortURL: "3",
				LongURL:  "https://github.com/",
				UserID:   "KS097f1lS&F",
			},
		},
		{
			name:      "Successfully added link",
			longURL:   "https://github.com/1234",
			userID:    "KS097f1lS&F",
			wantShort: "4",
			wantContain: model.URL{
				ShortURL: "4",
				LongURL:  "https://github.com/1234",
				UserID:   "KS097f1lS&F",
			},
		},
	}
	mock := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := mock.NewID()
			id62 := encoders.ToRBase62(id)
			url := model.URL{ShortURL: id62, LongURL: tt.longURL, UserID: tt.userID}
			_, _ = mock.Add(&url, nil)
			assert.Equal(t, tt.wantShort, url.ShortURL)
			var urls []model.URL
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
		want    string
		wantErr bool
	}{
		{
			name:    "Id in repo",
			short:   "2",
			want:    "https://yandex.ru/",
			wantErr: false,
		},
		{
			name:    "Id not in repo",
			short:   "3",
			want:    "",
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
			got, err := r.Get(tt.short, nil)
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
		want    []*model.URL
		wantErr bool
	}{
		{
			name:   "User has 2 urls",
			userID: "KS097f1lS&F",
			want: []*model.URL{
				{
					ShortURL: "1",
					LongURL:  "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders",
				},

				{
					ShortURL: "2",
					LongURL:  "https://yandex.ru/",
				},
			},
			wantErr: false,
		},
		{
			name:    "User has no urls",
			userID:  "SDADAD&FS()AS",
			want:    nil,
			wantErr: true,
		},
	}
	mock := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockRepo{
				urlsByShort: mock.urlsByShort,
				urlsByUser:  mock.urlsByUser,
			}
			got, err := r.GetByUser(tt.userID, nil)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestInMemRepo_NewID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "Test #1",
			want: 3,
		},
	}
	mock := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := mock.NewID()
			assert.Equal(t, tt.want, id)
		})
	}
}
