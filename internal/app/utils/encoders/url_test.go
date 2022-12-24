package encoders

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToBase62(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{name: "Test #1",
			url:  "yandex.ru",
			want: "SAAZrGBT2O5",
		},
		{name: "Test #2",
			url:  "github.com",
			want: "aAE3t8nGJ9A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToRBase62(tt.url))
		})
	}
}
