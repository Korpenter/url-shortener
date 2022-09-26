package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFileRepo(t *testing.T) {
	tests := []struct {
		name string
		urls []url
		want map[string]string
	}{
		{
			name: "Load preexisting file, No addition",
			want: map[string]string{"1": "yandex.ru", "2": "hero.ru"},
		},
		{
			name: " Add links and Load preexisting file.",
			urls: []url{
				{ID: "3", LongURL: "hell.ru"},
				{ID: "4", LongURL: "nvidia.ru"},
			},
			want: map[string]string{"1": "yandex.ru", "2": "hero.ru", "3": "hell.ru", "4": "nvidia.ru"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fRepo, err := newMockFileRepo()
			require.NoError(t, err)
			for _, v := range tt.urls {
				_, err = fRepo.add(v.LongURL, v.ID)
				require.NoError(t, err)
			}
			err = fRepo.load()
			require.NoError(t, err)
			assert.Equal(t, tt.want, fRepo.cache)
			err = fRepo.deleteMock()
			require.NoError(t, err)
		})
	}
}
