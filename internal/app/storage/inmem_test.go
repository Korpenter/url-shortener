package storage

import (
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TODO rewrite test with mocking
func TestInMemRepo_Add(t *testing.T) {
	type args struct {
		longURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Successfully added link",
			args: args{longURL: "https://github.com/"},
			want: "3",
		},
		{
			name: "Successfully added link",
			args: args{longURL: "https://github.com/"},
			want: "4",
		},
	}
	mockRepo := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := encoders.ToRBase62(mockRepo.NewID())
			assert.Equal(t, tt.want, mockRepo.Add(tt.args.longURL, id))
			assert.Contains(t, mockRepo.urls, tt.want)
		})
	}
}

func TestInMemRepo_Get(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Id in repo",
			args:    args{id: "2"},
			want:    "https://yandex.ru/",
			wantErr: false,
		},
		{
			name:    "Id not in repo",
			args:    args{id: "3"},
			want:    "",
			wantErr: true,
		},
	}
	mockRepo := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemRepo{
				urls: mockRepo.urls,
			}
			got, err := r.Get(tt.args.id)
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
	mockRepo := NewMockRepo()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mockRepo.NewID())
		})
	}
}
