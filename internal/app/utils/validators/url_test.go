package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Correct URL #1",
			url:  "https://habr.com/ru/post/541676/",
			want: true,
		},
		{
			name: "Correct URL #2",
			url:  "habr.com/ru/post/541676/",
			want: true,
		},
		{
			name: "Incorrect URL #1",
			url:  "habr",
			want: false,
		},
		{
			name: "Incorrect URL #2",
			url:  "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsURL(tt.url))
		})
	}
}
